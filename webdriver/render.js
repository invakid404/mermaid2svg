const [source, options, callback] = arguments;
window.mermaid.initialize(options);

(async () => {
  const { svg } = await window.mermaid.render("container", source);
  document.documentElement.innerHTML = svg;

  const container = document.querySelector("#container");
  const bbox = container.getBBox();
  container.setAttribute(
    "viewBox",
    [bbox.x, bbox.y, bbox.width, bbox.height].join(" ")
  );

  callback(container.outerHTML);
})();
