export function createMonthCursor(value = new Date()) {
  return new Date(Date.UTC(value.getFullYear(), value.getMonth(), 1));
}

export function shiftMonth(date, delta) {
  return new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth() + delta, 1));
}

export function formatMonthYear(date) {
  return new Intl.DateTimeFormat("en-US", {
    month: "long",
    year: "numeric",
    timeZone: "UTC",
  }).format(date);
}

export function formatCurrency(value) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(Number(value) || 0);
}

export function formatShortDate(value) {
  return new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  }).format(new Date(value));
}

export function toDateInputValue(value) {
  const date = value ? new Date(value) : new Date();
  return date.toISOString().slice(0, 10);
}

export function toAPIDate(value) {
  return new Date(`${value}T12:00:00.000Z`).toISOString();
}

export function isSameMonth(dateValue, cursor) {
  const date = new Date(dateValue);
  return (
    date.getUTCFullYear() === cursor.getUTCFullYear() &&
    date.getUTCMonth() === cursor.getUTCMonth()
  );
}
