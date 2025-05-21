import { register } from "./register.js";
import { login } from "./login.js";
import {
  registerForm,
  loginForm,
  postForm,
  postCard,
  Header,
  Footer,
  filterForm,
  activeUsersComponent,
  chatUsersComponent,

} from "./componnents.js";

import {
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
  getActiveUsers
} from "./api.js";

function styleBody() {
  // Clear the body
  document.body.innerHTML = "";
  document.body.style.background =
    "radial-gradient(circle at center, rgb(37, 42, 61), rgb(65, 47, 102))";
  document.body.style.fontFamily = "Segoe UI, sans-serif";
  // document.body.style.color = "#ffffff"; // fallback text color

  // Add welcome message
  const welcomeText = document.createElement("h1");
  welcomeText.textContent = "Welcome to Our Forum!";
  welcomeText.style.textAlign = "center";
  welcomeText.style.fontSize = "2.5rem";
  welcomeText.style.marginBottom = "1rem";
  welcomeText.style.padding = "1rem";
  welcomeText.style.color = "#ffd369"; // bright accent (yellow-gold)
  document.body.appendChild(welcomeText);

  const subText = document.createElement("p");
  subText.textContent = "Please log in to continue";
  subText.style.textAlign = "center";
  subText.style.fontSize = "1.2rem";
  subText.style.marginBottom = "2rem";
  subText.style.color = "#e0e0e0"; // soft gray-white
  document.body.appendChild(subText);
}

async function renderHomePage(data) {
  let user = await isAouth();
  if (!user) {
    styleBody();
    showLoginForm();
    return;
  }
  document.body.innerHTML = "";
  removeOldeForms();
  document.body.appendChild(Header(user));
  let main = document.createElement("main");
  main.id = "main"
  let activeUsers = document.createElement("section");
  activeUsers.setAttribute("id", "active_users")
  main.appendChild(activeUsers);
  if (!data) {
    data = await fetchPosts();
  }

  // let posts = await fetchPosts();

  let section = document.createElement("section");
  section.setAttribute("class", "container");
  section.setAttribute("id", "container");
  // fetch active users


  let posts = document.createElement("div");
  posts.setAttribute("class", "posts");
  if (!data) {
    posts.textContent = "No posts yet.";
  } else {
    data.forEach((post) => {
      posts.appendChild(postCard(post));
    });
  }

  section.appendChild(posts);
  main.appendChild(section);
  document.body.appendChild(main);
  document.body.appendChild(Footer());
  bindLoginBtn();
  bindRegisterbtn();
  showPostForm();
  if (user) {
    showProfile();
    bindfiletrBtn();
    const socket = await establishConnection();
    listenChatBtn(socket);
    logOut(socket)

  }


  postActions();
}
// this function diplay the login form
function bindLoginBtn() {
  const login_btn = document.getElementById("login_btn");
  const login_btn_1 = document.getElementById("login_btn_1");
  let result_btn = login_btn;
  if (login_btn_1) {
    result_btn = login_btn_1;
  }
  if (result_btn) {
    result_btn.addEventListener("click", (e) => {
      showLoginForm();
    });
  }
}
function showLoginForm(errors) {
  removeOldeForms();
  let form = loginForm(errors);
  form.classList.add("active");
  document.body.appendChild(form);

  bindRegisterbtn();
  login();
}

// this function dipay the registration form
function bindRegisterbtn() {
  const register_btn = document.getElementById("register_btn");
  if (register_btn) {
    register_btn.addEventListener("click", (e) => {
      showRegisterForm();
    });
  }
}
function showRegisterForm(errors = {}) {
  removeOldeForms();

  let form = registerForm(errors);
  form.classList.add("active");
  document.body.appendChild(form);

  bindLoginBtn();
  register();
}

// this function diplay the craete post form
function showPostForm(errors = {}, openImmediately = false) {
  const craete_post_btn = document.getElementById("craete_post_btn");
  if (!craete_post_btn && !openImmediately) return;

  const openForm = () => {
    const container = document.getElementById("container");
    container.classList.add("modal-active");

    // Remove any existing form first
    removeOldeForms();

    const form = postForm(errors);
    form.classList.add("active");
    document.body.appendChild(form);
    createPost();

    const close_btn = document.getElementById("close-form");
    close_btn.addEventListener("click", () => {
      form.remove();
      container.classList.remove("modal-active");
    });
  };

  // If we're calling this after a failed submit, open form immediately
  if (openImmediately) {
    openForm();
  } else {
    craete_post_btn.addEventListener("click", openForm);
  }
}

function removeOldeForms() {
  let allforms = document.querySelectorAll(".modal-overlay"); // add dot to select by class
  if (allforms.length > 0) {
    allforms.forEach((form) => {
      form.remove();
    });
  }
  document.getElementById("categoryFilterPanel")?.remove();
  document.getElementById("chat_users")?.remove();
}

function showMessage(message) {
  const popup = document.createElement("div");
  popup.setAttribute("id", "message_popup");
  popup.innerHTML = `<h2>${message}</h2>`;
  document.body.appendChild(popup);

  // Automatically hide after 3 seconds
  setTimeout(() => {
    popup.remove();
  }, 3000);
}

function showProfile() {
  let userProfile = document.getElementById("user-profile");
  if (userProfile) {
    userProfile.addEventListener("click", () => {
      let underProfile = document.getElementById("underProfile");
      underProfile.classList.toggle("hidden");
    });
  }
}

function showErrorPage(error) {
  document.body.innerHTML = `
    <div class="error-container">
      <h1 class="error-code">${error.code}</h1>
      <p class="error-message">${error.message}</p>
      <button class="back-home-btn" onclick="location.href='/'">Back Home</button>
    </div>
  `;
}

function renderComments(comments, postId, post) {
  // Remove existing panel if open
  document.getElementById("comment-panel")?.remove();

  // Create panel
  const panel = document.createElement("div");
  panel.id = "comment-panel";
  panel.className = "comment-panel";

  // Header with close button
  const header = document.createElement("div");
  header.className = "comment-header";

  const title = document.createElement("h3");
  title.textContent = `Comments for Post #${postId}`;
  header.appendChild(title);

  const closeBtn = document.createElement("button");
  closeBtn.textContent = "✖";
  closeBtn.className = "close-comment-panel";
  closeBtn.addEventListener("click", () => panel.remove());
  header.appendChild(closeBtn);

  panel.appendChild(header);

  // Comment list
  const list = document.createElement("div");
  list.className = "comment-container";

  if (comments.length === 0) {
    const empty = document.createElement("p");
    empty.textContent = "No comments yet.";
    empty.className = "no-comments";
    list.appendChild(empty);
  } else {
    comments.forEach((comment) => {
      const commentEl = document.createElement("div");
      commentEl.className = "comment-item";
      commentEl.innerHTML = `<strong>${comment.Creator}</strong>: ${comment.Content}`;
      list.appendChild(commentEl);
    });
  }

  panel.appendChild(list);
  post.appendChild(panel);
}

function postActions() {
  document.querySelectorAll(".post-card").forEach((postCard) => {
    const postId = postCard.querySelector("#post-id")?.textContent;

    // Like
    postCard.querySelector(".fa-thumbs-up")?.addEventListener("click", () => {
      const likeIcon = postCard.querySelector(".fa-thumbs-up");
      const countSpan = postCard.querySelector(".like-count");
      let currentCount = parseInt(countSpan.textContent, 10) || 0;

      const alreadyLiked = likeIcon.classList.contains("liked");

      if (alreadyLiked) {
        let unlike = reactToPost(postId, "dislike");
        if (unlike) {
          likeIcon.classList.remove("liked");
          countSpan.textContent = Math.max(0, currentCount - 1);
        }
      } else {
        let like = reactToPost(postId, "like");
        if (like) {
          likeIcon.classList.add("liked");
          countSpan.textContent = currentCount + 1;
        }
      }
    });

    // Dislike
    postCard.querySelector(".fa-thumbs-down")?.addEventListener("click", () => {
      const dislikeIcon = postCard.querySelector(".fa-thumbs-down");
      const countSpan = postCard.querySelector(".dislike-count");
      let currentCount = parseInt(countSpan.textContent, 10) || 0;

      const alreadyDisliked = dislikeIcon.classList.contains("disliked");

      if (alreadyDisliked) {
        let undo = reactToPost(postId, "dislike");
        if (undo) {
          dislikeIcon.classList.remove("disliked");
          countSpan.textContent = Math.max(0, currentCount - 1);
        }
      } else {
        let dislike = reactToPost(postId, "dislike");
        if (dislike) {
          dislikeIcon.classList.add("disliked");
          countSpan.textContent = currentCount + 1;
        }
      }
    });

    // Show comments
    postCard.querySelector(".fa-comment")?.addEventListener("click", () => {
      showComments(postId, postCard);
    });

    // Send comment
    const sendBtn = postCard.querySelector(".comment-button");
    const input = postCard.querySelector(".comment-input");

    sendBtn?.addEventListener("click", () => {
      const comment = input.value.trim();
      if (comment !== "") {
        input.value = "";
        let commented = sendPostCommen(postId, comment);
        if (commented) {
          let commentCount = postCard.querySelector(".comment-count");
          let currentCount = parseInt(commentCount.textContent, 10) || 0;
          commentCount.textContent = currentCount + 1;
        }
      }
    });
  });
}

function bindfiletrBtn() {
  const filter_btn = document.getElementById("filter_btn");
  if (filter_btn) {
    filter_btn.addEventListener("click", (e) => {
      showFilterForm();
    });
  }
}
function showFilterForm() {
  removeOldeForms();
  let form = filterForm();
  document.body.appendChild(form);
  // Handle close
  const closeBtn = form.querySelector(".close-filter-btn");
  closeBtn.addEventListener("click", () => {
    form.remove();
  });

  // Handle filter submit
  const submitBtn = document.getElementById("applyFilter");
  submitBtn.addEventListener("click", async (e) => {
    e.preventDefault();
    const checked = form.querySelectorAll("input[name='categories']:checked");
    const selectedCategories = Array.from(checked).map((input) => input.value);

    console.log("Selected categories:", selectedCategories);
    const posts = await fetchFilteredPosts(selectedCategories);
    renderHomePage(posts);

    form.remove(); // Optional: remove popup after applying
  });
}
function OneOffline(user) {
  let activeSection = document.getElementById('active_users')
  console.log(activeSection, "avtive section");
  if (activeSection) {
    if (user.type == "Online") {
      
      console.log(user.data, "user Data");
      
      const activuser = activeUsersComponent(user.data)
      console.log(activuser,"active user");
      if (activuser) {
        activeSection.appendChild(activuser)
      }
    } else if (user.type = "Offline") {
      let offuser = document.getElementById(`active-${user.data}`)
      if (offuser) {
        offuser.remove()
      }
    }
  }
}
const listenChatBtn = (socket) => {
  const chatBtn = document.getElementById("chat_btn");

  chatBtn.addEventListener("click", async () => {
    removeOldeForms();
    const activeUsers = await getActiveUsers();
    document.querySelector(".posts")?.classList.add("hidden");
    let container = document.querySelector(".container");
    container.appendChild(
      chatUsersComponent(activeUsers, showChatWindow, socket)
    );

    // Add close functionality
    document.getElementById("close_chat").addEventListener("click", () => {

      document.querySelector("#chat_users")?.remove();
      document.querySelector(".posts")?.classList.remove("hidden");
    });
  });
};

const showChatWindow = (container, user, socket) => {
  container.querySelector(".chat-users-list")?.classList.add("hidden");
  let chatContainer = container.querySelector(".chat-container");
  let chatWindow = document.getElementById("chat_window");
  if (!chatWindow) {
    chatWindow = document.createElement("div");
    chatWindow.setAttribute("id", "chat_window");
    chatWindow.innerHTML = `
      <div class="chat-header">
      <div id="chat_user_profile">
        <div class="avatar-wrapper">
        <img class="user-avatar" src="./assets/avatar.png" alt="Profile picture of ${user.nickname}" />
        <span class="status-dot online"></span>
      </div>
      <span user-id="${user.id}" class="user-nickname" id="${user.nickname}">${user.nickname}</span>
      </div>
      <button class="primary-btn" id="close_messages"><i class="fas fa-times"></i></button>
      </div>
      <div class="chat-messages"></div>
      <div class="message-actions">
      <input type="text" id="message" class="chat-input" placeholder="Type a message..." />
      <button  id="sent-message" class="sent-message primary-btn"><i class="fa-solid fa-paper-plane"></i></button>
      </div>
    `;
    chatContainer.appendChild(chatWindow);
    chat(chatContainer, socket);
    let close = document.getElementById("close_messages");

    close.addEventListener("click", () => {

      chatWindow.remove();

      container.querySelector(".chat-users-list")?.classList.remove("hidden");
    });
  }
};
const removetyping = (container) => {
  let typers = container.querySelectorAll(".typing-indicator")
  console.log(typers, 'typers');

  if (typers.length > 0) {
    typers.forEach((elem) => {
      elem.remove()
    })
  }
}



export {
  renderHomePage,
  showLoginForm,
  showRegisterForm,
  showPostForm,
  bindRegisterbtn,
  bindLoginBtn,
  bindfiletrBtn,
  showMessage,
  showErrorPage,
  postActions,
  OneOffline,
  renderComments,
  removetyping
};
