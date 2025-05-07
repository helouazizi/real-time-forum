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
  showPostForm({},false);
});
