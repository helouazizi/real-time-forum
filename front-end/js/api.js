import {
  renderHomePage,
  showErrorPage,
  showLoginForm,
  showMessage,
  showPostForm,
  renderComments,
} from "./dom.js";
async function isAouth() {
  try {
    const response = await fetch("http://localhost:3000/api/v1/users/info", {
      method: "GET",
      credentials: "include", // 👈 This tells the browser to send cookies
    });
    if (response.ok) {
      let data = await response.json();
      console.log(data);

      sessionStorage.setItem("user_id", data.id);
      return data;
    } else {
      return null;
    }
  } catch (err) {
    console.log(err);
  }
}

function logOut() {
  let log_out_btn = document.getElementById("log_out");
  if (log_out_btn) {
    log_out_btn.addEventListener("click", async () => {
      try {
        const response = await fetch(
          "http://localhost:3000/api/v1/users/logout",
          {
            method: "GET",
            credentials: "include",
          }
        );

        if (response.ok) {
          renderHomePage();
        } else {
          const errorData = await response.json();
          throw { code: errorData.Code, message: errorData.Message };
        }
      } catch (err) {
        showErrorPage(err);
      }
    });
  }
}

function createPost() {
  const formElement = document.querySelector("[createPost_form_element]");
  if (!formElement) return;

  formElement.addEventListener("submit", async (e) => {
    e.preventDefault();
    const title = formElement.querySelector("#title").value.trim();
    const content = formElement.querySelector("#content").value.trim();
    const categoryCheckboxes = formElement.querySelectorAll(
      'input[name="categories"]:checked'
    );
    const categories = Array.from(categoryCheckboxes).map((cb) => cb.value);

    const postData = {
      title,
      content,
      categories,
    };
    try {
      const response = await fetch(
        "http://localhost:3000/api/v1/posts/create",
        {
          method: "POST",
          credentials: "include", // Very important
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(postData),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        if (errorData.UserErrors.HasError) {
          showPostForm(errorData.UserErrors, true);
          return;
        }
        if (errorData.Code === 401) {
          showLoginForm();
          return;
        }
        const error = {
          code: errorData.Code,
          message: errorData.Message,
        };
        throw error;
      }
      const result = await response.json();
      showMessage(result.Message);
      setTimeout(() => {
        renderHomePage();
      }, 2000);
    } catch (err) {
      showErrorPage(err);
    }
  });
}

async function fetchPosts() {
  try {
    const response = await fetch("http://localhost:3000/");
    if (!response.ok) {
      let err = {
        code: response.status,
        message: response.statusText,
      };
      throw err;
    }

    const posts = await response.json();
    return posts;
  } catch (error) {
    showErrorPage(error);
  }
}

async function fetchFilteredPosts(categories) {
  try {
    const response = await fetch("http://localhost:3000/api/v1/posts/filter", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ categories }), // send selected categories
    });

    if (!response.ok) {
      const err = {
        code: response.status,
        message: response.statusText,
      };
      throw err;
    }

    const posts = await response.json();
    return posts;
  } catch (error) {
    showErrorPage(error);
  }
}

async function reactToPost(postId, reaction) {
  try {
    const response = await fetch("http://localhost:3000/api/v1/posts/react", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        post_id: parseInt(postId),
        reaction: reaction,
      }),
    });

    if (!response.ok) {
      const errData = await response.json();
      let err = {
        code: errData.Code,
        message: errData.Message,
      };
      throw err;
    }
    return true;
  } catch (error) {
    showErrorPage(error);
  }
}
async function sendPostCommen(postId, commenttext) {
  console.log(postId, commenttext, "hhhh");

  try {
    const response = await fetch(
      "http://localhost:3000/api/v1/posts/addComment",
      {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          post_id: parseInt(postId),
          comment: commenttext,
        }),
      }
    );

    if (!response.ok) {
      const errData = await response.json();
      let err = {
        code: errData.Code,
        message: errData.Message,
      };
      throw err;
    }
    let res = await response.json();
    showMessage(res.Message);

    console.log("Comment submitted successfully:");
    // Optionally update the UI here
  } catch (error) {
    showErrorPage(error);
  }
}

async function showComments(postId, container) {
  try {
    const response = await fetch(
      `http://localhost:3000/api/v1/posts/fetchComments?postId=${postId}`,
      { credentials: "include" }
    );

    if (!response.ok) {
      const errData = await response.json();
      throw { code: errData.Code, message: errData.Message };
    }

    const comments = await response.json();
    if (!comments) return;
    console.log(comments);

    renderComments(comments, postId, container);
  } catch (error) {
    showErrorPage(error);
  }
}

async function chat() {
  const senderId = sessionStorage.getItem("user_id");

  const receiverId = 456;
  const messageInput = document.getElementById("message-input");

  // WebSocket connection
  const socket = new WebSocket("ws://http://localhost:3000/api/v1/chat");

  // When the WebSocket is open, send the message
  socket.onopen = function (event) {
    const sendButton = document.getElementById("send-button");
    sendButton.addEventListener("click", function () {
      // Get the message content
      const message = messageInput.value;
      const messageData = {
        sender_id: senderId,
        receiver_id: receiverId,
        message: message,
      };

      // Send the message
      socket.send(JSON.stringify(messageData));

      // Clear the input field after sending the message
      messageInput.value = "";
    });
  };

  // Handling incoming messages
  socket.onmessage = function (event) {
    const incomingMessage = JSON.parse(event.data);
    console.log("Received message:", incomingMessage);
  };
}

async function getActiveUsers() {
  try {
    const response = await fetch("http://localhost:3000/api/v1/active", {
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error("Failed to fetch active users");
    }

    const users = await response.json();
    console.log("Active users:", users);
    return users;
  } catch (error) {
    console.error("Error fetching active users:", error);
    return [];
  }
}

export {
  isAouth,
  logOut,
  createPost,
  fetchPosts,
  reactToPost,
  sendPostCommen,
  showComments,
  fetchFilteredPosts,
  getActiveUsers,
};
