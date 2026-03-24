import { createExpense, deleteExpense, listExpenses, login, register, updateExpense } from "./api.js";
import {
  createMonthCursor,
  formatCurrency,
  formatMonthYear,
  formatShortDate,
  isSameMonth,
  shiftMonth,
  toAPIDate,
  toDateInputValue,
} from "./date.js";
import { clearSession, getSession, setSession } from "./store.js";
import { animateSwap, setLoading, showToast } from "./ui.js";

const state = {
  authMode: "login",
  route: "dashboard",
  expenses: [],
  monthCursor: createMonthCursor(),
};

const CATEGORIES = [
  "Food & Dining",
  "Groceries",
  "Transportation",
  "Entertainment",
  "Shopping",
  "Bills & Utilities",
  "Health",
  "Travel",
  "Education",
  "Personal Care",
  "Gifts",
  "Home",
  "Subscriptions",
  "Other",
];

const elements = {
  authView: document.querySelector("#auth-view"),
  appView: document.querySelector("#app-view"),
  authForm: document.querySelector("#auth-form"),
  authSubmit: document.querySelector("#auth-submit"),
  authCopy: document.querySelector("#auth-copy"),
  loginTab: document.querySelector("#login-tab"),
  registerTab: document.querySelector("#register-tab"),
  navLinks: Array.from(document.querySelectorAll(".nav-link[data-route]")),
  dashboardView: document.querySelector("#dashboard-view"),
  expensesView: document.querySelector("#expenses-view"),
  monthLabel: document.querySelector("#month-label"),
  monthlyTotal: document.querySelector("#monthly-total"),
  monthlyCount: document.querySelector("#monthly-count"),
  averageExpense: document.querySelector("#average-expense"),
  topCategory: document.querySelector("#top-category"),
  largestExpense: document.querySelector("#largest-expense"),
  largestExpenseCaption: document.querySelector("#largest-expense-caption"),
  recentActivity: document.querySelector("#recent-activity"),
  recentActivityCaption: document.querySelector("#recent-activity-caption"),
  categoryBreakdown: document.querySelector("#category-breakdown"),
  recentExpenses: document.querySelector("#recent-expenses"),
  allExpensesCount: document.querySelector("#all-expenses-count"),
  tableBody: document.querySelector("#expense-table-body"),
  prevMonth: document.querySelector("#prev-month"),
  nextMonth: document.querySelector("#next-month"),
  logoutButton: document.querySelector("#logout-button"),
  openCreateModal: document.querySelector("#open-create-modal"),
  openModalButtons: Array.from(document.querySelectorAll("[data-open-modal]")),
  modal: document.querySelector("#expense-modal"),
  modalTitle: document.querySelector("#modal-title"),
  modalClose: document.querySelector("#modal-close"),
  modalCancel: document.querySelector("#modal-cancel"),
  expenseForm: document.querySelector("#expense-form"),
  expenseSubmit: document.querySelector("#expense-submit"),
  expenseCategory: document.querySelector("#expense-category"),
};

function normalizeRoute(route) {
  return route === "expenses" ? "expenses" : "dashboard";
}

function populateCategoryOptions() {
  elements.expenseCategory.innerHTML = [
    '<option value="" disabled selected>Select a category</option>',
    ...CATEGORIES.map((category) => `<option value="${escapeAttribute(category)}">${escapeHTML(category)}</option>`),
  ].join("");
}

function getFilteredExpenses() {
  return state.expenses.filter((expense) => isSameMonth(expense.date, state.monthCursor));
}

function setAuthMode(mode) {
  state.authMode = mode;
  const isLogin = mode === "login";
  elements.loginTab.classList.toggle("is-active", isLogin);
  elements.registerTab.classList.toggle("is-active", !isLogin);
  elements.authSubmit.textContent = isLogin ? "Login" : "Create account";
  elements.authCopy.textContent = isLogin
    ? "Sign in to continue to your monthly dashboard."
    : "Create an account to start tracking expenses.";
  document.querySelector("#password").setAttribute("autocomplete", isLogin ? "current-password" : "new-password");
}

function setRoute(route) {
  const nextRoute = normalizeRoute(route);
  state.route = nextRoute;
  elements.navLinks.forEach((link) => {
    link.classList.toggle("is-active", link.dataset.route === nextRoute);
  });
  elements.dashboardView.classList.toggle("hidden", nextRoute !== "dashboard");
  elements.expensesView.classList.toggle("hidden", nextRoute !== "expenses");
  window.location.hash = nextRoute;
}

function getTopCategory(expenses) {
  const totals = new Map();
  expenses.forEach((expense) => {
    totals.set(expense.category, (totals.get(expense.category) || 0) + Number(expense.amount));
  });

  let topName = "None";
  let topTotal = 0;

  totals.forEach((total, category) => {
    if (total > topTotal) {
      topName = category;
      topTotal = total;
    }
  });

  return topName;
}

function getLargestExpense(expenses) {
  return expenses.reduce((largest, expense) => {
    if (!largest || Number(expense.amount) > Number(largest.amount)) {
      return expense;
    }
    return largest;
  }, null);
}

function getRecentActivity(expenses) {
  if (!expenses.length) {
    return { label: "Quiet", caption: "No recent expenses yet." };
  }

  const latest = expenses.reduce((current, expense) => {
    return new Date(expense.date) > new Date(current.date) ? expense : current;
  }, expenses[0]);

  return {
    label: formatShortDate(latest.date),
    caption: `${latest.vendor} in ${latest.category}`,
  };
}

function renderCategoryBreakdown(expenses) {
  const container = elements.categoryBreakdown;
  if (!expenses.length) {
    container.className = "category-breakdown empty-state";
    container.textContent = "No expenses yet for this month.";
    return;
  }

  const totals = new Map();
  expenses.forEach((expense) => {
    totals.set(expense.category, (totals.get(expense.category) || 0) + Number(expense.amount));
  });

  const categories = Array.from(totals.entries()).sort((a, b) => b[1] - a[1]);
  const grandTotal = categories.reduce((sum, [, amount]) => sum + amount, 0);

  container.className = "category-breakdown";
  container.innerHTML = categories
    .map(([category, amount]) => {
      const width = grandTotal ? Math.max(8, Math.round((amount / grandTotal) * 100)) : 0;
      return `
        <article class="category-row">
          <div class="category-row-header">
            <strong>${escapeHTML(category)}</strong>
            <span>${formatCurrency(amount)}</span>
          </div>
          <div class="progress-track">
            <div class="progress-fill" style="width: ${width}%"></div>
          </div>
        </article>
      `;
    })
    .join("");
}

function renderRecentExpenses(expenses) {
  const container = elements.recentExpenses;
  if (!expenses.length) {
    container.className = "recent-list empty-state";
    container.textContent = "Add your first expense to start tracking.";
    return;
  }

  container.className = "recent-list";
  container.innerHTML = expenses
    .slice(0, 5)
    .map(
      (expense) => `
        <article class="recent-card">
          <div class="recent-card-header">
            <strong>${escapeHTML(expense.vendor)}</strong>
            <span>${formatCurrency(expense.amount)}</span>
          </div>
          <span>${escapeHTML(expense.category)} · ${formatShortDate(expense.date)}</span>
          <span class="muted-copy">${escapeHTML(expense.description || "No description")}</span>
        </article>
      `,
    )
    .join("");
}

function renderExpenseTable() {
  const rows = state.expenses;
  elements.allExpensesCount.textContent = `${rows.length} ${rows.length === 1 ? "transaction" : "transactions"}`;
  if (!rows.length) {
    elements.tableBody.innerHTML = `
      <tr>
        <td colspan="6" class="empty-state">No expenses yet.</td>
      </tr>
    `;
    return;
  }

  elements.tableBody.innerHTML = rows
    .map(
      (expense) => `
        <tr>
          <td data-label="Date">${formatShortDate(expense.date)}</td>
          <td data-label="Vendor">${escapeHTML(expense.vendor)}</td>
          <td data-label="Category">${escapeHTML(expense.category)}</td>
          <td data-label="Amount">${formatCurrency(expense.amount)}</td>
          <td data-label="Description">${escapeHTML(expense.description || "—")}</td>
          <td class="actions-column" data-label="Actions">
            <div class="expense-row-actions">
              <button class="chip-button" data-modal-edit="${expense.id}" type="button">Edit</button>
              <button class="chip-button danger" data-delete-expense="${expense.id}" type="button">Delete</button>
            </div>
          </td>
        </tr>
      `,
    )
    .join("");
}

function renderDashboard() {
  const expenses = getFilteredExpenses();
  const total = expenses.reduce((sum, expense) => sum + Number(expense.amount), 0);
  const average = expenses.length ? total / expenses.length : 0;
  const largestExpense = getLargestExpense(expenses);
  const recentActivity = getRecentActivity(expenses);

  elements.monthLabel.textContent = formatMonthYear(state.monthCursor);
  elements.monthlyTotal.textContent = formatCurrency(total);
  elements.monthlyCount.textContent = String(expenses.length);
  elements.averageExpense.textContent = formatCurrency(average);
  elements.topCategory.textContent = getTopCategory(expenses);
  elements.largestExpense.textContent = formatCurrency(largestExpense?.amount || 0);
  elements.largestExpenseCaption.textContent = largestExpense
    ? `${largestExpense.vendor} on ${formatShortDate(largestExpense.date)}`
    : "No expenses this month.";
  elements.recentActivity.textContent = recentActivity.label;
  elements.recentActivityCaption.textContent = recentActivity.caption;
  renderCategoryBreakdown(expenses);
  renderRecentExpenses(expenses);
  renderExpenseTable();
  animateSwap(elements.dashboardView);
}

function renderApp() {
  renderDashboard();
}

function openExpenseModal(expense = null) {
  elements.modalTitle.textContent = expense ? "Update expense" : "Add expense";
  elements.expenseSubmit.textContent = expense ? "Save changes" : "Save expense";
  elements.expenseForm.dataset.expenseId = expense ? String(expense.id) : "";
  elements.expenseForm.vendor.value = expense?.vendor || "";
  elements.expenseCategory.value = expense?.category || "";
  elements.expenseForm.amount.value = expense?.amount || "";
  elements.expenseForm.date.value = toDateInputValue(expense?.date);
  elements.expenseForm.description.value = expense?.description || "";
  elements.modal.showModal();
}

function closeExpenseModal() {
  elements.modal.close();
  elements.expenseForm.reset();
  elements.expenseForm.dataset.expenseId = "";
  elements.expenseCategory.value = "";
}

async function refreshExpenses() {
  const expenses = await listExpenses();
  state.expenses = expenses;
  renderApp();
}

async function handleAuthSubmit(event) {
  event.preventDefault();
  setLoading(elements.authForm, true);

  const formData = new FormData(elements.authForm);
  const credentials = {
    username: String(formData.get("username") || "").trim(),
    password: String(formData.get("password") || ""),
  };

  try {
    const response = state.authMode === "login" ? await login(credentials) : await register(credentials);
    setSession({ username: credentials.username, jwtToken: response.jwtToken });
    elements.authForm.reset();
    showToast(state.authMode === "login" ? "Logged in successfully." : "Account created.");
    await bootAuthenticatedApp();
  } catch (error) {
    showToast(error.message, "error");
  } finally {
    setLoading(elements.authForm, false);
  }
}

async function handleExpenseSubmit(event) {
  event.preventDefault();
  setLoading(elements.expenseForm, true);

  const formData = new FormData(elements.expenseForm);
  const payload = {
    vendor: String(formData.get("vendor") || "").trim(),
    category: String(formData.get("category") || "").trim(),
    amount: Number(formData.get("amount") || 0),
    date: toAPIDate(String(formData.get("date") || "")),
    description: String(formData.get("description") || "").trim(),
  };

  const expenseId = elements.expenseForm.dataset.expenseId;

  try {
    if (expenseId) {
      await updateExpense(expenseId, payload);
      showToast("Expense updated.");
    } else {
      await createExpense(payload);
      showToast("Expense added.");
    }
    closeExpenseModal();
    await refreshExpenses();
  } catch (error) {
    showToast(error.message, "error");
  } finally {
    setLoading(elements.expenseForm, false);
  }
}

async function handleTableClick(event) {
  const modalEditId = event.target.closest("[data-modal-edit]")?.dataset.modalEdit;
  if (modalEditId) {
    const expense = state.expenses.find((entry) => String(entry.id) === modalEditId);
    if (expense) {
      openExpenseModal(expense);
    }
    return;
  }
  const deleteId = event.target.closest("[data-delete-expense]")?.dataset.deleteExpense;
  if (deleteId) {
    await handleDeleteExpense(Number(deleteId));
  }
}

async function handleDeleteExpense(expenseId) {
  const confirmed = window.confirm("Delete this expense?");
  if (!confirmed) {
    return;
  }

  try {
    await deleteExpense(expenseId);
    showToast("Expense deleted.");
    await refreshExpenses();
  } catch (error) {
    showToast(error.message, "error");
  }
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function escapeAttribute(value) {
  return escapeHTML(value);
}

function showAuthView() {
  elements.authView.classList.remove("hidden");
  elements.appView.classList.add("hidden");
}

function showAppView() {
  elements.authView.classList.add("hidden");
  elements.appView.classList.remove("hidden");
}

async function bootAuthenticatedApp() {
  const session = getSession();
  if (!session?.jwtToken) {
    showAuthView();
    return;
  }

  showAppView();
  try {
    await refreshExpenses();
    setRoute(normalizeRoute(window.location.hash.replace("#", "")));
  } catch (error) {
    clearSession();
    showAuthView();
    showToast(error.message || "Session expired.", "error");
  }
}

function bindEvents() {
  elements.loginTab.addEventListener("click", () => setAuthMode("login"));
  elements.registerTab.addEventListener("click", () => setAuthMode("register"));
  elements.authForm.addEventListener("submit", handleAuthSubmit);
  elements.prevMonth.addEventListener("click", () => {
    state.monthCursor = shiftMonth(state.monthCursor, -1);
    renderDashboard();
  });
  elements.nextMonth.addEventListener("click", () => {
    state.monthCursor = shiftMonth(state.monthCursor, 1);
    renderDashboard();
  });
  elements.navLinks.forEach((link) => {
    link.addEventListener("click", () => setRoute(link.dataset.route));
  });
  elements.logoutButton.addEventListener("click", () => {
    clearSession();
    state.expenses = [];
    showAuthView();
    setAuthMode("login");
    window.location.hash = "";
  });
  elements.openCreateModal.addEventListener("click", () => openExpenseModal());
  elements.openModalButtons.forEach((button) => {
    button.addEventListener("click", () => openExpenseModal());
  });
  elements.modalClose.addEventListener("click", closeExpenseModal);
  elements.modalCancel.addEventListener("click", closeExpenseModal);
  elements.expenseForm.addEventListener("submit", handleExpenseSubmit);
  elements.tableBody.addEventListener("click", handleTableClick);
  window.addEventListener("hashchange", () => {
    const nextRoute = normalizeRoute(window.location.hash.replace("#", ""));
    if (getSession()?.jwtToken) {
      setRoute(nextRoute);
    }
  });
}

function init() {
  populateCategoryOptions();
  setAuthMode("login");
  bindEvents();
  bootAuthenticatedApp();
}

init();
