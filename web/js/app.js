(function () {
  "use strict";

  // ---------- API helper ----------
  const API_BASE = ""; // same-origin, served by the Go server itself

  function getToken() {
    return localStorage.getItem("ts_token");
  }

  function setSession(token, email) {
    localStorage.setItem("ts_token", token);
    if (email) localStorage.setItem("ts_email", email);
  }

  function clearSession() {
    localStorage.removeItem("ts_token");
    localStorage.removeItem("ts_email");
  }

  async function api(path, options) {
    options = options || {};
    const headers = Object.assign({ "Content-Type": "application/json" }, options.headers || {});
    const token = getToken();
    if (token) headers["Authorization"] = "Bearer " + token;

    const res = await fetch(API_BASE + path, {
      method: options.method || "GET",
      headers,
      body: options.body ? JSON.stringify(options.body) : undefined,
    });

    let data = null;
    try {
      data = await res.json();
    } catch (e) {
      data = null;
    }

    if (!res.ok) {
      const message = (data && data.error) || "Request failed (" + res.status + ")";
      const err = new Error(message);
      err.status = res.status;
      throw err;
    }
    return data;
  }

  // ---------- DOM refs ----------
  const authView = document.getElementById("auth-view");
  const appView = document.getElementById("app-view");

  const tabBtns = document.querySelectorAll(".tab-btn");
  const loginForm = document.getElementById("login-form");
  const registerForm = document.getElementById("register-form");
  const loginError = document.getElementById("login-error");
  const registerError = document.getElementById("register-error");

  const userEmailEl = document.getElementById("user-email");
  const logoutBtn = document.getElementById("logout-btn");

  const ticketListEl = document.getElementById("ticket-list");
  const ticketCountEl = document.getElementById("ticket-count");
  const listErrorEl = document.getElementById("list-error");
  const emptyStateEl = document.getElementById("empty-state");
  const filterChips = document.querySelectorAll(".filter-chip");

  const newTicketBtn = document.getElementById("new-ticket-btn");
  const emptyNewBtn = document.getElementById("empty-new-btn");
  const newTicketModal = document.getElementById("new-ticket-modal");
  const newTicketForm = document.getElementById("new-ticket-form");
  const newTicketError = document.getElementById("new-ticket-error");

  const detailModal = document.getElementById("detail-modal");
  const detailTitle = document.getElementById("detail-title");
  const detailStatus = document.getElementById("detail-status");
  const detailId = document.getElementById("detail-id");
  const detailCreated = document.getElementById("detail-created");
  const detailUpdated = document.getElementById("detail-updated");
  const detailDescription = document.getElementById("detail-description");
  const detailError = document.getElementById("detail-error");
  const statusSetBtns = document.querySelectorAll(".status-set-btn");

  const toastEl = document.getElementById("toast");

  let currentTickets = [];
  let currentFilter = "all";
  let activeTicketId = null;

  // ---------- Toast ----------
  let toastTimer = null;
  function showToast(message, type) {
    toastEl.textContent = message;
    toastEl.className = "toast" + (type ? " " + type : "");
    toastEl.classList.remove("hidden");
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toastEl.classList.add("hidden"), 3000);
  }

  // ---------- View switching ----------
  function showApp() {
    authView.classList.add("hidden");
    appView.classList.remove("hidden");
    userEmailEl.textContent = localStorage.getItem("ts_email") || "";
    loadTickets();
  }

  function showAuth() {
    appView.classList.add("hidden");
    authView.classList.remove("hidden");
  }

  // ---------- Auth tabs ----------
  tabBtns.forEach((btn) => {
    btn.addEventListener("click", () => {
      tabBtns.forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");
      const tab = btn.dataset.tab;
      loginForm.classList.toggle("hidden", tab !== "login");
      registerForm.classList.toggle("hidden", tab !== "register");
      loginError.textContent = "";
      registerError.textContent = "";
    });
  });

  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    loginError.textContent = "";
    const email = document.getElementById("login-email").value.trim();
    const password = document.getElementById("login-password").value;
    try {
      const data = await api("/auth/login", { method: "POST", body: { email, password } });
      setSession(data.token, email);
      showApp();
    } catch (err) {
      loginError.textContent = err.message;
    }
  });

  registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    registerError.textContent = "";
    const email = document.getElementById("register-email").value.trim();
    const password = document.getElementById("register-password").value;
    try {
      await api("/auth/register", { method: "POST", body: { email, password } });
      // Auto-login after successful registration
      const data = await api("/auth/login", { method: "POST", body: { email, password } });
      setSession(data.token, email);
      showApp();
    } catch (err) {
      registerError.textContent = err.message;
    }
  });

  logoutBtn.addEventListener("click", () => {
    clearSession();
    showAuth();
  });

  // ---------- Tickets ----------
  function formatDate(iso) {
    if (!iso) return "—";
    const d = new Date(iso);
    if (isNaN(d.getTime())) return iso;
    return d.toLocaleString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  function statusLabel(status) {
    return String(status || "").replace("_", " ");
  }

  async function loadTickets() {
    listErrorEl.textContent = "";
    try {
      const tickets = await api("/tickets");
      currentTickets = Array.isArray(tickets) ? tickets : [];
      renderTicketList();
    } catch (err) {
      if (err.status === 401) {
        clearSession();
        showAuth();
        return;
      }
      listErrorEl.textContent = err.message;
    }
  }

  function renderTicketList() {
    const filtered =
      currentFilter === "all" ? currentTickets : currentTickets.filter((t) => t.status === currentFilter);

    ticketCountEl.textContent = currentTickets.length + (currentTickets.length === 1 ? " ticket" : " tickets");

    ticketListEl.innerHTML = "";

    if (filtered.length === 0) {
      emptyStateEl.classList.remove("hidden");
      ticketListEl.classList.add("hidden");
      return;
    }
    emptyStateEl.classList.add("hidden");
    ticketListEl.classList.remove("hidden");

    // Newest first
    const sorted = [...filtered].sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

    sorted.forEach((ticket) => {
      const card = document.createElement("div");
      card.className = "ticket-card";
      card.innerHTML =
        '<div class="ticket-card-main">' +
        '<p class="ticket-card-title"></p>' +
        '<p class="ticket-card-meta"></p>' +
        "</div>" +
        '<span class="status-badge"></span>';

      card.querySelector(".ticket-card-title").textContent = ticket.title;
      card.querySelector(".ticket-card-meta").textContent =
        "#" + ticket.id + " · Created " + formatDate(ticket.created_at);

      const badge = card.querySelector(".status-badge");
      badge.textContent = statusLabel(ticket.status);
      badge.classList.add(ticket.status);

      card.addEventListener("click", () => openDetail(ticket.id));
      ticketListEl.appendChild(card);
    });
  }

  filterChips.forEach((chip) => {
    chip.addEventListener("click", () => {
      filterChips.forEach((c) => c.classList.remove("active"));
      chip.classList.add("active");
      currentFilter = chip.dataset.status;
      renderTicketList();
    });
  });

  // ---------- New ticket modal ----------
  function openModal(el) {
    el.classList.remove("hidden");
  }
  function closeModal(el) {
    el.classList.add("hidden");
  }

  document.querySelectorAll("[data-close]").forEach((btn) => {
    btn.addEventListener("click", () => closeModal(document.getElementById(btn.dataset.close)));
  });
  [newTicketModal, detailModal].forEach((overlay) => {
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) closeModal(overlay);
    });
  });

  newTicketBtn.addEventListener("click", () => {
    newTicketForm.reset();
    newTicketError.textContent = "";
    openModal(newTicketModal);
  });
  emptyNewBtn.addEventListener("click", () => {
    newTicketForm.reset();
    newTicketError.textContent = "";
    openModal(newTicketModal);
  });

  newTicketForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    newTicketError.textContent = "";
    const title = document.getElementById("new-ticket-title").value.trim();
    const description = document.getElementById("new-ticket-desc").value.trim();
    try {
      await api("/tickets", { method: "POST", body: { title, description } });
      closeModal(newTicketModal);
      showToast("Ticket created", "success");
      loadTickets();
    } catch (err) {
      newTicketError.textContent = err.message;
    }
  });

  // ---------- Detail modal ----------
  async function openDetail(id) {
    activeTicketId = id;
    detailError.textContent = "";
    detailTitle.textContent = "Loading…";
    openModal(detailModal);
    try {
      const ticket = await api("/tickets/" + id);
      renderDetail(ticket);
    } catch (err) {
      detailError.textContent = err.message;
    }
  }

  function renderDetail(ticket) {
    detailTitle.textContent = ticket.title;
    detailStatus.textContent = statusLabel(ticket.status);
    detailStatus.className = "status-badge " + ticket.status;
    detailId.textContent = "#" + ticket.id;
    detailCreated.textContent = formatDate(ticket.created_at);
    detailUpdated.textContent = formatDate(ticket.updated_at);
    detailDescription.textContent = ticket.description || "No description provided.";

    statusSetBtns.forEach((btn) => {
      btn.classList.toggle("current", btn.dataset.status === ticket.status);
    });
  }

  statusSetBtns.forEach((btn) => {
    btn.addEventListener("click", async () => {
      if (!activeTicketId || btn.classList.contains("current")) return;
      detailError.textContent = "";
      try {
        const ticket = await api("/tickets/" + activeTicketId + "/status", {
          method: "PATCH",
          body: { status: btn.dataset.status },
        });
        renderDetail(ticket);
        showToast("Status updated", "success");
        loadTickets();
      } catch (err) {
        detailError.textContent = err.message;
      }
    });
  });

  // ---------- Init ----------
  if (getToken()) {
    showApp();
  } else {
    showAuth();
  }
})();
