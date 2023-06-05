const [source, options, callback] = arguments;
window.mermaid.initialize(options);

(async () => {
  const { svg } = await window.mermaid.render("container", source);
  document.documentElement.innerHTML = svg;

  const container = document.querySelector("#container");

  container.removeAttribute("style");

  const bbox = container.getBBox({
    fill: true,
    stroke: true,
    markers: true,
    clipped: true,
  });
  const rect = container.getBoundingClientRect();

  container.setAttribute(
    "viewBox",
    [
      Math.floor(bbox.x) - 5,
      Math.floor(bbox.y) - 5,
      Math.ceil(Math.min(bbox.width, rect.width)) + 10,
      Math.ceil(bbox.height) + 10,
    ].join(" ")
  );

  callback(container.outerHTML);
})();
