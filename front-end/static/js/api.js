import {
  renderHomePage,
  showErrorPage,
  showMessage,
  showPostForm,
  renderComments,
  removetyping,
  OneOffline
} from "./dom.js";
import { createTypingIndicator } from "./componnents.js";

async function isAouth() {
  try {
    const response = await fetch("http://localhost:3000/api/v1/users/info", {
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
      console.log(response.status, "status");


      if (!response.ok) {
        const errorData = await response.json();
        console.log(errorData, "errrrrr");

        if (errorData.UserErrors.HasError) {
          showPostForm(errorData.UserErrors, true);
          return;
        }
        // if (errorData.Code === 401) {
        //   showLoginForm();
        //   return;
        // }
        const error = {
          code: errorData.Code,
          message: errorData.Message,
        };
        console.log(error, "test");

        throw error;
      }
      const result = await response.json();
      showMessage(result.Message);
      setTimeout(() => {
        renderHomePage();
      }, 2000);
    } catch (err) {
      console.log(err, "ctacj");

      showErrorPage(err);
    }
  });
}

async function fetchPosts() {
  try {
    const response = await fetch("http://localhost:3000/api/v1/posts");
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
    renderComments(comments, postId, container);
  } catch (error) {
    showErrorPage(error);
  }
}

async function chat(chatContainer, socket) {
  const username = document.getElementById("username").innerText
  const senderId = parseInt(sessionStorage.getItem("user_id"));
  const spanElement = chatContainer.querySelector("span[user-id]");
  const receiverId = parseInt(spanElement.getAttribute("user-id"));

  // Pagination state
  let offset = 0;
  const limit = 10;
  let loading = false;

  const messagesContainer = chatContainer.querySelector(".chat-messages");
  // Load chat history (initial fetch)
  function fetchHistory() {
    const historyRequest = {
      username: username,
      sender: senderId,
      receiver: receiverId,
      message: "",
      limit,
      offset,
    };
    socket.send(JSON.stringify(historyRequest));
  }

  fetchHistory(); // Load the latest messages on open

  // Scroll listener for loading older messages
  messagesContainer.addEventListener("scroll", () => {
    if (messagesContainer.scrollTop === 0 && !loading) {
      loading = true;
      fetchHistory();
    }
  });
  setupTypingIndicator(socket, username, senderId, receiverId)
  // Send message button logic
  const sendBtn = chatContainer.querySelector("#sent-message");
  sendBtn.addEventListener("click", () => {
    const messageInput = chatContainer.querySelector("#message");
    const message = messageInput.value.trim();

    if (!message) return;

    const messageData = {
      data: {
        username: username,
        sender: senderId,
        receiver: receiverId,
        message,
        timestamp: new Date(Date.now()).toLocaleString(),

      },
      type: 'message'

    };

    socket.send(JSON.stringify(messageData.data));

    // Display message immediately in UI
    const messageElement = createMessageElement(messageData, senderId);
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
    messageInput.value = "";
  });

  // Handle incoming WebSocket messages
  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    const chatWindowExists = chatContainer.querySelector("#chat_window");
    const type = data.type;
    console.log(data.data);
    

    if (data.data.typing == "typing") {
    
    let typingcontainer = createTypingIndicator(data.data.username)
        removetyping(messagesContainer)
        messagesContainer.appendChild(typingcontainer)
    } else {
        removetyping(messagesContainer)
    }

    if (type === "history") {

      loading = false;
      const currentScrollHeight = messagesContainer.scrollHeight;
      const msgEl = createMessageElement(data, senderId);
      messagesContainer.prepend(msgEl);

      messagesContainer.scrollTop =
        messagesContainer.scrollHeight - currentScrollHeight;

      offset += 1;



    } else if (type === "message" && !data.data.typing) {

        removetyping(messagesContainer)
      const msgEl = createMessageElement(data, senderId);
      if (chatWindowExists) {
        messagesContainer.appendChild(msgEl);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      } else {
        showMessage(
          `New message from ${data.data.username}`
        );
      }
    } else if(type== "Online" || type== "Offline" ){
      OneOffline(data,)
    }
  };

  socket.onerror = (error) => {
    console.error("WebSocket error:", error);
  };

  socket.onclose = () => {
    removetyping(messagesContainer)
    console.log("WebSocket connection closed");
  };
}
function createMessageElement(data, senderId) {
  const messageElement = document.createElement("div");
  let senderName;
  if (data.type == "history") {

    messageElement.className = data.data.sender === parseInt(senderId) ? "outgoing-message" : "incoming-message";


    senderName = data.data.sender === parseInt(senderId) ? "You" : data.data.RecieverNickname || "Unknown";


  } else {

    messageElement.className = data.data.sender === parseInt(senderId) ? "outgoing-message" : "incoming-message";

    senderName = data.data.sender === parseInt(senderId) ? "You" : data.data.username || "Unknown";
  }
  const timestamp = new Date(data.data.timestamp).toLocaleString(); // Assumes ISO timestamp
  messageElement.innerHTML = `
    <div class="message-meta">
      <strong class="message-user">${senderName}</strong>
      <time class="message-time">${timestamp}</time>
    </div>
    <div class="message-content">
      <p>${sanitizeHTML(data.data.message)}</p>
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
  const username = document.getElementById("username").innerText
  const senderId = sessionStorage.getItem("user_id");
  const socket = new WebSocket("ws://localhost:3000/api/v1/chat");

  return new Promise((resolve, reject) => {
    socket.onopen = () => {
      socket.send(JSON.stringify({ sender: parseInt(senderId), username: username }));
      resolve(socket);
    };
    socket.onmessage = (event) => {
      const ResponseMessages = JSON.parse(event.data);
      if (ResponseMessages.data.message) {
        showMessage(
          `New message from ${ResponseMessages.data.username}`
        );
      }
    };

    socket.onerror = reject;
  });
}

function setupTypingIndicator(socket, username, senderId, receiverId) {
  const messageInput = document.getElementById('message');
  if (!messageInput) return;

  let typingTimeout;
  let isTyping = false; setupTypingIndicator

  messageInput.addEventListener("input", () => {
    if (!isTyping) {
      // Send typing: true only once when typing starts
      const typingData = {
        username,
        sender: senderId,
        receiver: receiverId,
        typing: "typing"
      };
      socket.send(JSON.stringify(typingData));
      isTyping = true;
    }

    // Reset the timer on every keypress
    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
      // Send typing: false after 2 seconds of inactivity
      const stopTypingData = {
        username,
        sender: senderId,
        receiver: receiverId,
        typing: "fin"
      };
      socket.send(JSON.stringify(stopTypingData));
      isTyping = false;
    }, 2000); // 2 seconds delay
  });
  window.addEventListener("beforeunload", ()=> {
    const stopTypingData = {
        username,
        sender: senderId,
        receiver: receiverId,
        typing: "fin"
      };
      socket.send(JSON.stringify(stopTypingData));
      isTyping = false;
  })
 

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

  chat,
  establishConnection,
};
