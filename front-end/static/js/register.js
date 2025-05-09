import {
  showRegisterForm,
  showLoginForm,
  showMessage,
  showErrorPage,
} from "./dom.js";

function register() {
  const registerFormElement = document.getElementById("register_form_element");
  registerFormElement.addEventListener("submit", async (e) => {
    e.preventDefault(); // Stop regular form submission

    const formData = new FormData(registerFormElement);
    const data = Object.fromEntries(formData.entries());
    data.age = parseInt(data.age);

    try {
      const response = await fetch(
        "http://localhost:3000/api/v1/users/register",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        if (errorData.UserErrors.HasError) {
          showRegisterForm(errorData.UserErrors);
          return;
        }
        const error = {
          code: errorData.Code,
          message: errorData.Message,
        };
        throw error;
      }
      const dataa = await response.json();
      setTimeout(() => {
        showMessage(dataa.Message);
      }, 2000);
      setTimeout(() => {
        showLoginForm();
      }, 1000);
    } catch (err) {
      showErrorPage(err);
    }
  });
}

export { register };
