const Header = (user) => {
  let header = document.createElement("header");

  if (user) {
    header.innerHTML = `
      <h1 class="logo"><a href="/">Forum</a></h1>
      <div class="header-center-buttons">
       <button class="primary-btn new_post_btn" id="craete_post_btn">
       <i class="fas fa-plus"></i>
       </button>
        <button class="primary-btn  filter_btn" id="filter_btn">
          <i class="fas fa-filter"></i>
        </button>
        <button class="primary-btn  chat_btn" id="chat_btn">
          <i class="fas fa-comments"></i>
        </button>
      </div>
      <nav class="navigation-links">
        <div class="user-profile" id="user-profile">
          <img src="/front-end/static/assets/avatar.png" alt="User Profile" class="profile-pic" />
        </div>
      </nav>
    `;
    header.appendChild(userProfile(user));
  } else {
    header.innerHTML = `
      <h1 class="logo"><a href="/front-end/">Forum</a></h1>
      <nav class="navigation-links">
        <button class="primary-btn" id="login_btn">
          <i class="fas fa-sign-in-alt"></i> Login
        </button>
      </nav>
    `;
  }

  return header;
};

const userProfile = (user) => {
  let underProfile = document.createElement("div");
  underProfile.setAttribute("class", "underProfile hidden");
  underProfile.setAttribute("id", "underProfile");

  underProfile.innerHTML = `
    <div class="profile-card">
      <div class="profile-header">
        <img src="/front-end/static/assets/avatar.png" alt="User Profile" class="profile-pic" />
        <div>
          <h2 id="username">${user.nickname}</h2>
          <p>${user.email}</p>
        </div>
      </div>

      <div class="profile-details">
        <div class="detail-item"><strong>First Name:</strong> ${user.first_name
    }</div>
        <div class="detail-item"><strong>Last Name:</strong> ${user.last_name
    }</div>
        <div class="detail-item"><strong>Gender:</strong> ${user.gender}</div>
        <div class="detail-item"><strong>Age:</strong> ${user.age}</div>
        <div class="detail-item"><strong>Created At:</strong> ${new Date(
      user.created_at
    ).toLocaleString()}</div>
        <div class="detail-item"><strong>Updated At:</strong> ${new Date(
      user.updated_at
    ).toLocaleString()}</div>
      </div>

      <div class="profile-actions">
        <button class="primary-btn" id="log_out">Log Out</button>
      </div>
    </div>
  `;

  return underProfile;
};

const loginForm = (errrors = {}) => {
  let form = document.createElement("div");
  form.setAttribute("class", "modal-overlay");
  form.setAttribute("id", "login_form");
  form.innerHTML = `
                <div class="modal-content">
                    <!-- Login Form -->
                <form action="http://localhost:3000/api/v1/users/login" method="POST" id="login_form_element">
                <h2><i class="fas fa-sign-in-alt"></i> Login</h2>
    
                <label for="login_id">Nickname or E-mail</label>
                <input type="text" id="login_id" name="login_id" required />
                <span>${errrors.Nickname ? errrors.Nickname : ""}</span>
    
                <label for="login_password">Password</label>
                <input type="password" id="login_password" name="password" required />
                <span>${errrors.Pass ? errrors.Pass : ""}</span>
    
                <button type="submit">Login</button>
                </form>
                <div class="register_action"> 
                    <p>Don't have an account ? </p>
                    <button class="primary-btn" id="register_btn"><i class="fas fa-user-plus"></i> register</button >
                </div>
                </div>
            `;
  return form;
};

const registerForm = (errrors = {}) => {
  let form = document.createElement("div");
  form.setAttribute("class", "modal-overlay");
  form.setAttribute("id", "register_form");
  form.innerHTML = `
                <div class="modal-content">
                    <!-- Login Form -->
                <form action="http://localhost:3000/api/v1/users/register" method="POST" id="register_form_element">
                <h2><i class="fas fa-user-plus"></i> Register</h2>
                
                <label for="nickname">Nickname</label>
                <input type="text" id="nickname" name="nickname" required>
                <span>${errrors.Nickname ? errrors.Nickname : ""}</span>

                <label for="age">Age</label>
                <input type="number" id="age" name="age" min="1" max="120" required>
                <span>${errrors.Age ? errrors.Age : ""}</span>

                <label for="gender">Gender</label>
                <select id="gender" name="gender" required>
                <option value="" disabled selected>Select gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                </select>
                <span>${errrors.Gender ? errrors.Gender : ""}</span>

      
                <label for="first_name">First Name</label>
                <input type="text" id="first_name" name="first_name" required>
                <span>${errrors.LastName ? errrors.FirstName : ""}</span>
      
                <label for="last_name">Last Name</label>
                <input type="text" id="last_name" name="last_name" required>
                <span>${errrors.LastName ? errrors.LastName : ""}</span>
      
                <label for="email">E-mail</label>
                <input type="email" id="email" name="email" required>
                <span>${errrors.Email ? errrors.Email : ""}</span>
      
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required>
                <span>${errrors.Pass ? errrors.Pass : ""}</span>
      
                <button type="submit">Register</button>
                </form>
                <div class="register_action"> 
                    <p>Alredy have an acoount ? </p>
                    <button class="primary-btn" id="login_btn_1"><i class="fas fa-sign-in-alt"></i> login</button >
                </div>
                </div>
            `;

  return form;
};

const postCard = (post) => {
  let postElement = document.createElement("div");
  postElement.setAttribute("class", "post-card");

  // Format the date
  const date = new Date(post.CreatedAt).toLocaleDateString();

  // Convert categories array to HTML span elements
  const categoryTags = post.categories
    .map((cat) => `<span class="category-tag">${cat}</span>`)
    .join("");

  postElement.innerHTML = `
    <div class="post-header">
      <img
        src="/front-end/static/assets/avatar.png"
        alt="User Profile"
        class="profile-pic"
      />
      <div class="user-info">
        <h4 class="username">${post.Creator}</h4>
        <span class="post-date">Posted on ${date}</span>
        <span  class="hidden" id="post-id">${post.ID}</span>
      </div>
    </div>

    <div class="post-body">
      <h3 class="post-title">${post.title}</h3>
      <p class="post-content">
        ${post.content}
      </p>
    </div>

    <div class="post-categories">
      ${categoryTags}
    </div>

    <div class="comment-section">
      <input type="text" placeholder="Write a comment..." class="comment-input" />
      <button class="comment-button"><i class="fa-solid fa-paper-plane"></i></button>
    </div>

    <div class="post-footer">
      <span><i class="fa-solid fa-comment"></i> <span class="comment-count">${post.TotalComments}</span></span>
    </div>
  `;

  return postElement;
};

const postForm = (errors = {}) => {
  // let careate our form

  let form = document.createElement("div");
  form.setAttribute("class", "modal-overlay");
  form.setAttribute("id", "post_form");
  form.innerHTML = `
       <div class="modal-content">
         <button class="close-btn" id="close-form">&times;</button>
         <form action="/create_post" method="POST" createPost_form_element>
           <h2><i class="fas fa-plus"> </i>New Post</h2>
       
           <label for="title">Title:</label>
           <input type="text" id="title" name="title" maxlength="255" required />
           <span>${errors.PostTilte ? errors.PostTilte : ""}</span>

         
           <label for="content">Content:</label>
           <textarea id="content" name="content" rows="6" required></textarea>
           <span>${errors.PostContent ? errors.PostContent : ""}</span>
         
           <label>Select Categories:</label>
           <div class="category-container">
           <div class="category-checkbox">
               <input type="checkbox" id="cat-tech" name="categories" value="Technology" />
               <label for="cat-tech">Technology</label>
           </div>
           <div class="category-checkbox">
               <input type="checkbox" id="cat-sci" name="categories" value="Science" />
               <label for="cat-sci">Science</label>
           </div>
           <div class="category-checkbox">
               <input type="checkbox" id="cat-health" name="categories" value="Health" />
               <label for="cat-health">Health</label>
           </div>
           <div class="category-checkbox">
               <input type="checkbox" id="cat-life" name="categories" value="Lifestyle" />
               <label for="cat-life">Lifestyle</label>
           </div>
           <div class="category-checkbox">
               <input type="checkbox" id="cat-edu" name="categories" value="Education" />
               <label for="cat-edu">Education</label>
           </div>
           <div class="category-checkbox">
               <input type="checkbox" id="cat-game" name="categories" value="Gaming" />
               <label for="cat-game">Gaming</label>
           </div>
           <div class="category-checkbox">
               <input type="checkbox" id="cat-biz" name="categories" value="Business" />
               <label for="cat-biz">Business</label>
           </div>
           </div>
           <span>${errors.Postcategories ? errors.Postcategories : ""}</span>
   
         
           <button type="submit">Create Post</button>
         </form>
         
       </div>
           `;
  return form;
};
const activeUsersComponent = (users) => {
  const container = document.createElement("div");
  container.setAttribute("id", "active_users");

  const usersHTML = users
    .map((user) => {
      if (user.is_active) {

        return `
          <li class="user-item">
            <div class="avatar-wrapper">
              <img class="user-avatar" src="/front-end/static/assets/avatar.png" alt="Profile picture of ${user.nickname}" />
              ${user.is_active ? '<span class="status-dot active"></span>' : ''}
            </div>
            <span class="user-nickname">${user.nickname}</span>
          </li>
        `;
      }
    })
    .join("");

  container.innerHTML = `
    <ul class="users-list">${usersHTML}</ul>
  `;

  return container;
};

const chatUsersComponent = (users, onUserClick, socket) => {
  const container = document.createElement("div");
  container.setAttribute("id", "chat_users");
  container.setAttribute("class", "chat_users");
  const usersHTML = users
    .map(
      (user) => `
    <li class="chat-user-item" data-user="${user.nickname}">
      <div class="avatar-wrapper">
        <img class="user-avatar" src="/front-end/static/assets/avatar.png" alt="Profile picture of ${user.nickname}" />
        ${user.is_active ? '<span class="status-dot active"></span>' : ''}
      </div>
      <span  class="user-nickname">${user.nickname}</span>
    </li>
  `
    )
    .join("");

  container.innerHTML = `
    <div class="chat-container"><ul class="chat-users-list">${usersHTML}</ul></div>
    <button class="primary-btn" id="close_chat"><i class="fas fa-times"></i> Close Chat</button>

  `;

  // Event delegation for click on any user
  container.addEventListener("click", (e) => {
    const item = e.target.closest(".chat-user-item");
    if (item) {
      const nickname = item.dataset.user;
      const selectedUser = users.find((u) => u.nickname === nickname);
      if (selectedUser) {
        onUserClick(container, selectedUser, socket); // open chat window
      }
    }
  });

  return container;
};


const filterForm = () => {
  const container = document.createElement("form");
  container.id = "categoryFilterPanel";
  container.className = "category-filter-form";
  container.setAttribute("aria-label", "Filter posts by category");

  container.innerHTML = `
    <div class="filter-header">
      <h2>Filter Posts</h2>
    </div>

    <div class="form-group">
      <label for="category-list">Select Categories:</label>
      <div class="category-container" id="category-list">
        ${[
      { id: "tech", name: "Technology" },
      { id: "sci", name: "Science" },
      { id: "health", name: "Health" },
      { id: "life", name: "Lifestyle" },
      { id: "edu", name: "Education" },
      { id: "game", name: "Gaming" },
      { id: "biz", name: "Business" },
    ]
      .map(
        (cat) => `
          <div class="category-checkbox">
            <input type="checkbox" id="cat-${cat.id}" name="categories" value="${cat.name}" />
            <label for="cat-${cat.id}">${cat.name}</label>
          </div>`
      )
      .join("")}
      </div>
    </div>

    <div class="form-buttons">
      <button type="submit" id="applyFilter" class="primary-btn">Apply Filter</button>
      <button type="button" id="closeFilter" class="close-filter-btn">Cancel</button>
    </div>
  `;

  return container;
};

const Footer = () => {
  let footer = document.createElement("footer");
  footer.innerHTML = `
    <p>&copy; 2025 Forum. All rights reserved.</p>
    `;
  return footer;
};
function createTypingIndicator(username) {

  const container = document.createElement("div");
  container.classList.add("typing-indicator");
  container.classList.add("incoming-message");
  let user = document.createElement("span");
  user.innerText = username + "  is typing"
  container.appendChild(user)
  for (let i = 0; i < 3; i++) {
    const dot = document.createElement("span");
    dot.className = `dot dot-${i + 1}`;
    container.appendChild(dot);
  }
 

  return container;
}


export {
  createTypingIndicator,
  Footer,
  Header,
  loginForm,
  registerForm,
  postCard,
  postForm,
  filterForm,
  activeUsersComponent,
  chatUsersComponent
};
