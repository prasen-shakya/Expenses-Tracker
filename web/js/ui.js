export function showToast(message, variant = "success") {
  const region = document.querySelector("#toast-region");
  const toast = document.createElement("div");
  toast.className = `toast ${variant === "error" ? "error" : ""}`.trim();
  toast.textContent = message;
  region.appendChild(toast);

  window.setTimeout(() => {
    toast.remove();
  }, 2800);
}

export function setLoading(element, isLoading) {
  element.classList.toggle("loading", isLoading);
}

export function animateSwap(element) {
  element.classList.remove("fade-slide");
  void element.offsetWidth;
  element.classList.add("fade-slide");
}
