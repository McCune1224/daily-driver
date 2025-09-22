document.addEventListener("DOMContentLoaded", function() {
  let dots = 0;
  const maxDots = 3;
  let loadingInterval;

  function startLoadingAnimation() {
    const loadingDots = document.getElementById("loading-dots");
    loadingInterval = setInterval(() => {
      dots = (dots + 1) % (maxDots + 1);
      loadingDots.textContent = ".".repeat(dots);
    }, 100);
  }

  function stopLoadingAnimation() {
    clearInterval(loadingInterval);
    const loadingDots = document.getElementById("loading-dots");
    if (loadingDots) loadingDots.textContent = "";
  }

  document.body.addEventListener("htmx:indicator", function(evt) {
    const indicator = document.getElementById("loading-indicator");
    if (evt.detail.visible) {
      if (indicator) indicator.style.display = "";
      startLoadingAnimation();
    } else {
      if (indicator) indicator.style.display = "none";
      stopLoadingAnimation();
    }
  });
});

