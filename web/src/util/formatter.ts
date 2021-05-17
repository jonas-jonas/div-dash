export function formatMoney(amount: number) {
  return new Intl.NumberFormat("de-DE", {
    style: "currency",
    currency: "EUR",
  }).format(amount);
}

export function formatDate(date: string) {
  const d = Date.parse(date);
  return new Intl.DateTimeFormat("de-DE", {
    dateStyle: "short",
  }).format(d);
}

export function formatTime(date: string) {
  const d = Date.parse(date);
  return new Intl.DateTimeFormat("de-DE", {
    timeStyle: "medium",
  }).format(d);
}

export function formatPercent(percent: number) {
  return new Intl.NumberFormat("de-DE", {
    style: "percent",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(percent);
}
