import {
  renderHomePage,
  showErrorPage,
  showMessage,
  showPostForm,
  renderComments,
  removetyping,
  OneOffline,
  showChatWindow,
  renderprofiles,
} from "./dom.js";
import { createTypingIndicator, chatUsersComponent } from "./componnents.js";

async function isAouth() {
  try {
    const response = await fetch("/api/v1/users/info", {
      method: "GET",
      credentials: "include", // 👈 This tells the browser to send cookies
    });
    if (response.ok) {
      let data = await response.json();
      sessionStorage.setItem("user_id", data.id);
      return data;
    } else {
      return null;
    }
  } catch (err) {
    console.log(err, "rrrrrrrrrr");

    showErrorPage(err);
  }
}

function logOut(socket) {
  let log_out_btn = document.getElementById("log_out");
  if (log_out_btn) {
    log_out_btn.addEventListener("click", async () => {
      try {
        const response = await fetch("/api/v1/users/logout", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          renderHomePage();
          socket.close();
          localStorage.removeItem("is_logged");
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
      const response = await fetch("/api/v1/posts/create", {
        method: "POST",
        credentials: "include", // Very important
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(postData),
      });

      if (!response.ok) {
        const errorData = await response.json();

        if (errorData.UserErrors.HasError) {
          showPostForm(errorData.UserErrors, true);
          return;
        }

        const error = {
          code: errorData.Code,
          message: errorData.Message,
        };

        throw error;
      }
      const result = await response.json();

      renderHomePage();
    } catch (err) {
      showErrorPage(err);
    }
  });
}

async function fetchPosts() {
  try {
    const response = await fetch("/api/v1/posts");
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
    const response = await fetch("/api/v1/posts/filter", {
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
async function getActiveUsers() {
  const senderId = sessionStorage.getItem("user_id");
  try {
    const response = await fetch("/api/v1/users", {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ sender_id: parseInt(senderId) }),
    });

    if (!response.ok) {
      throw { code: response.Code, message: response.Messgae };
    }

    const users = await response.json();
    return users;
  } catch (error) {
    showErrorPage(error);
    return [];
  }
}

async function reactToPost(postId, reaction) {
  try {
    const response = await fetch("/api/v1/posts/react", {
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
  try {
    const response = await fetch("/api/v1/posts/addComment", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        post_id: parseInt(postId),
        comment: commenttext,
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
    let res = await response.json();
    showMessage(res.Message);
    // Optionally update the UI here
  } catch (error) {
    showErrorPage(error);
  }
}

async function showComments(postId, container) {
  try {
    const response = await fetch(
      `/api/v1/posts/fetchComments?postId=${postId}`,
      { credentials: "include" }
    );

    if (!response.ok) {
      const errData = await response.json();
      throw { code: errData.Code, message: errData.Message };
    }

    const comments = await response.json();
    if (!comments) return;
    renderComments(comments, postId, container);
  } catch (error) {
    showErrorPage(error);
  }
}
// Debounce helper
function debounce(func, delay) {
  let timeout;
  return function (...args) {
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), delay);
  };
}

async function chat(chatContainer, socket) {
  const username = document.getElementById("username").innerText;
  const senderId = parseInt(sessionStorage.getItem("user_id"));
  const spanElement = chatContainer.querySelector("span[user-id]");
  const receiverId = parseInt(spanElement.getAttribute("user-id"));

  let offset = 0;
  const limit = 10;
  let loading = false;
  const messagesContainer = chatContainer.querySelector(".chat-messages");

  async function fetchHistory() {
    if (loading) return;
    loading = true;
    try {
      const requestBody = {
        sender_id: senderId,
        receiver_id: receiverId,
        offset: offset,
        limit: limit,
      };

      const response = await fetch("/api/v1/chat/history", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(requestBody),
      });

      const dataa = await response.json();

      if (Array.isArray(dataa) && dataa.length > 0) {
        const currentScrollHeight = messagesContainer.scrollHeight;

        dataa.reverse().forEach((data) => {
          const messageElement = createMessageElement(data, senderId);
          messagesContainer.prepend(messageElement);
        });

        messagesContainer.scrollTop =
          messagesContainer.scrollHeight - currentScrollHeight;
        offset += limit;
      }
    } catch (err) {
      console.error("Error fetching history:", err);
    }

    loading = false;
  }

  // Debounced scroll handler
  messagesContainer.addEventListener(
    "scroll",
    debounce(() => {
      if (messagesContainer.scrollTop === 0) {
        fetchHistory();
      }
    }, 500)
  );

  // Initial fetch
  fetchHistory();

  // Typing indicator
  setupTypingIndicator(socket, username, senderId, receiverId);

  // Send message logic
  const sendBtn = chatContainer.querySelector("#sent-message");
  sendBtn.addEventListener("click", () => {
    const messageInput = chatContainer.querySelector("#message");
    const message = messageInput.value.trim();
    if (!message) return;

    const messageData = {
      sender_id: senderId,
      receiver_id: receiverId,
      message_content: message,
      timestamp: new Date(Date.now()).toLocaleString(), // can cause an eror
      message_type: "message",
    };
    socket.send(JSON.stringify(messageData));
    offset++;
    messageInput.value = "";
  });

  // Handle WebSocket messages
  socket.onmessage = async (event) => {
    const message = JSON.parse(event.data);
    const type = message.message_type;
    console.log(message, "message chat ");
    const chatWindowExists = chatContainer.querySelector("#chat_window");

   

    if (type === "typing") {
      const typingContainer = createTypingIndicator(message.sender_nickname);
      removetyping(messagesContainer);
      messagesContainer.appendChild(typingContainer);
    } else {
      removetyping(messagesContainer);
    }

    if (type === "message" /*&& !data.data.typing*/) {
      const msgEl = createMessageElement(message, senderId);
      if (chatWindowExists) {
        messagesContainer.appendChild(msgEl);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      } else {
        showMessage(`New message from ${message.sender_nickname}`);
        renderprofiles(socket, type);
      }
    }

    if (type === "Online" || type === "Offline") {
      OneOffline(message);

      renderprofiles(socket, type);
    }
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  socket.onclose = () => {
    removetyping(messagesContainer);

    console.log("WebSocket connection closed");
  };
  document.getElementById("close_chat").addEventListener("click", () => {
    document.querySelector("#chat_users")?.remove();
    document.querySelector(".posts")?.classList.remove("hidden");
  });
}

function createMessageElement(data, senderId) {
  const messageElement = document.createElement("div");
  let senderName;
  if (data.type == "history") {
    messageElement.className =
      data.sender_id === parseInt(senderId)
        ? "outgoing-message"
        : "incoming-message";
    senderName =
      data.sender_id === parseInt(senderId) ? "You" : data.reciever_nickname;
  } else {
    messageElement.className =
      data.sender_id === parseInt(senderId)
        ? "outgoing-message"
        : "incoming-message";
    senderName =
      data.sender_id === parseInt(senderId)
        ? "You"
        : data.sender_nickname || "Unknown";
  }

  const timestamp = new Date(data.timestamp).toLocaleString(); // Assumes ISO timestamp
  messageElement.innerHTML = `
    <div class="message-meta">
      <strong class="message-user">${senderName}</strong>
      <time class="message-time">${timestamp}</time>
    </div>
    <div class="message-content">
      <p>${sanitizeHTML(data.message_content)}</p>
    </div>
  `;

  return messageElement;
}

// Sanitize to prevent XSS
function sanitizeHTML(str) {
  const div = document.createElement("div");
  div.textContent = str;
  return div.innerHTML;
}

async function establishConnection() {
  const senderId = sessionStorage.getItem("user_id");
  const socket = new WebSocket("/api/v1/chat");

  return new Promise((resolve, reject) => {
    socket.onopen = () => {
      socket.send(JSON.stringify({ sender_id: parseInt(senderId) }));
      resolve(socket);
    };
    socket.onmessage = async (event) => {
      const ResponseMessages = JSON.parse(event.data);
      if (
        ResponseMessages.message_type == "Online" ||
        ResponseMessages.message_type == "Offline"
      ) {
        OneOffline(ResponseMessages);
        renderprofiles(socket, ResponseMessages.message_type);
      }
      if (ResponseMessages.message_content) {
        showMessage(`New message from ${ResponseMessages.sender_nickname}`);
        let chatusers = document.getElementById("chat_users");

        if (chatusers) {
          const activeUsers = await getActiveUsers();
          console.log(activeUsers, "activeUsers");

          let container = chatUsersComponent(
            activeUsers,
            showChatWindow,
            socket
          );
          console.log(container, "container");
          chatusers.replaceWith(container);
        }
      }
    };
    socket.onerror = reject;
  });
}

function setupTypingIndicator(socket, username, senderId, receiverId) {
  const messageInput = document.getElementById("message");
  if (!messageInput) return;

  let typingTimeout;
  let isTyping = false;
  setupTypingIndicator;

  messageInput.addEventListener("input", () => {
    if (!isTyping) {
      // Send typing: true only once when typing starts
      const typingData = {
        sender_nickname: username,
        sender_id: senderId,
        receiver_id: receiverId,
        message_type: "typing",
      };
      socket.send(JSON.stringify(typingData));
      isTyping = true;
    }

    // Reset the timer on every keypress
    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
      // Send typing: false after 2 seconds of inactivity
      const stopTypingData = {
        sender_nickname: username,
        sender_id: senderId,
        receiver_id: receiverId,
        message_type: "fin",
      };
      socket.send(JSON.stringify(stopTypingData));
      isTyping = false;
    }, 2000); // 2 seconds delay
  });
  window.addEventListener("beforeunload", () => {
    const stopTypingData = {
      sender_nickname: username,
      sender_id: senderId,
      receiver_id: receiverId,
      message_type: "fin",
    };
    socket.send(JSON.stringify(stopTypingData));
    isTyping = false;
  });
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
  chat,
  establishConnection,
};
