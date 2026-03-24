import { getSession } from "./store.js";

async function request(path, options = {}) {
  const session = getSession();
  const headers = new Headers(options.headers || {});

  if (!headers.has("Content-Type") && options.body) {
    headers.set("Content-Type", "application/json");
  }

  if (session?.jwtToken) {
    headers.set("Authorization", `Bearer ${session.jwtToken}`);
  }

  const response = await fetch(path, {
    ...options,
    headers,
  });

  const isJSON = response.headers.get("content-type")?.includes("application/json");
  const payload = isJSON ? await response.json() : null;

  if (!response.ok) {
    const message = payload?.error || payload?.message || "Request failed";
    throw new Error(message);
  }

  return payload;
}

export async function login(credentials) {
  return request("/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  });
}

export async function register(credentials) {
  return request("/register", {
    method: "POST",
    body: JSON.stringify(credentials),
  });
}

export async function listExpenses() {
  return request("/expenses", { method: "GET" });
}

export async function createExpense(expense) {
  return request("/expenses", {
    method: "POST",
    body: JSON.stringify(expense),
  });
}

export async function updateExpense(id, expense) {
  return request(`/expenses/${id}`, {
    method: "PUT",
    body: JSON.stringify(expense),
  });
}

export async function deleteExpense(id) {
  return request(`/expenses/${id}`, {
    method: "DELETE",
  });
}
