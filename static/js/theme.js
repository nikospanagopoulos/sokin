(function () {
  var themes = [
    { id: "light-sand", name: "Light Sand" },
    { id: "dark-sea", name: "Dark Sea" },
    { id: "mountain", name: "Mountain" },
  ];
  var KEY = "sokin-theme";
  var root = document.documentElement;
  var current = themes[0];

  function setLabel() {
    var el = document.getElementById("theme-name");
    if (el) el.textContent = current.name;
  }

  function apply(theme) {
    current = theme;
    root.setAttribute("data-theme", theme.id);
    setLabel();
  }

  var saved = null;
  try { saved = localStorage.getItem(KEY); } catch (e) {}
  for (var i = 0; i < themes.length; i++) {
    if (themes[i].id === saved) { current = themes[i]; break; }
  }
  root.setAttribute("data-theme", current.id);

  document.addEventListener("DOMContentLoaded", function () {
    setLabel();
    var btn = document.getElementById("theme-toggle");
    if (!btn) return;
    btn.addEventListener("click", function () {
      var next = themes[(themes.indexOf(current) + 1) % themes.length];
      apply(next);
      try { localStorage.setItem(KEY, next.id); } catch (e) {}
    });
  });
})();