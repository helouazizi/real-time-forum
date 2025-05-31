import {
  renderHomePage,
  bindLoginBtn,
  bindRegisterbtn,
  showPostForm,
} from "./dom.js";

document.addEventListener("DOMContentLoaded", async () => {
  await renderHomePage();
  bindLoginBtn();
  bindRegisterbtn();
  showPostForm({}, false);
  window.addEventListener("storage", async (e) => {
    if (e.key === "is_logged") {
      if (e.newValue === null) {
        await renderHomePage();
      }
    }
  });
});
