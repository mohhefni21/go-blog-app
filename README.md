
<br />
<div align="center">
  <h3 align="center">Blog Api</h3>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

This is a Blog API built using Golang with a Domain-Driven Design (DDD) approach. The API provides a robust platform for users to authenticate, create blog posts, add comments, interact with posts through likes, shares, and bookmarks, and categorize posts with tags.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Features
1. Authentication (Auth)
    * Register
    * Login
    * Refresh token
    * Login with google
    * Logout
    * Update profile
2. Posts
    * Add post
    * Get all post
    * Get detail post
    * Delete post
    * Update post
    * Upload content Image (CKEditor api spesification)
3. Comments
    * Add comment
    * Update comment
    * Delete comment
4. Interactions
    * Add Interaction(like, share and bookmard)
    * Delete Interaction
5. Tags
    * Get tags

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

This project is based on the following packages:

* Golang
* Echo(golang framework)
* Golang-migrate
* Google Oauth
* JWT
* Postgres

<p align="right">(<a href="#readme-top">back to top</a>)</p>
