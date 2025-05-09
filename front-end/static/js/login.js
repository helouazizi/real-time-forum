import { showLoginForm, showErrorPage, renderHomePage } from "./dom.js";
function login() {
  const loginFormElement = document.getElementById("login_form_element");
  loginFormElement.addEventListener("submit", async (e) => {
    e.preventDefault(); // Stop regular form submission

    const formData = new FormData(loginFormElement);
    const data = Object.fromEntries(formData.entries());
    try {
      const response = await fetch("http://localhost:3000/api/v1/users/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Very important
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        const errorData = await response.json();
        if (errorData.UserErrors.HasError) {
          showLoginForm(errorData.UserErrors);
          return;
        }
        const error = {
          code: errorData.Code,
          message: errorData.Message,
        };
        throw error;
      }
      location.reload()
      renderHomePage()
    } catch (err) {
      showErrorPage(err);
    }
  });
}

export { login };
