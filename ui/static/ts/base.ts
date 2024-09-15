import { submit_form, validate_form } from "./form";
import { InputErrorElement } from "./error";

const signupFormEl: HTMLFormElement | any =
  document.querySelector("#signup form");
const signupBtnEl: HTMLButtonElement | any = document.querySelector(
  "#signup button[type=submit]"
);
const loginFormEl: HTMLFormElement | any =
  document.querySelector("#login form");
const loginBtnEl: HTMLButtonElement | any = document.querySelector(
  "#login button[type=submit]"
);
const bookFormEl: HTMLFormElement | any =
  document.querySelector("#book-create form");
const bookBtnEl: HTMLButtonElement | any = document.querySelector(
  "#book-create button[type=submit]"
);

if (window.location.href == "http://localhost:8080/users/signup") {
  validate_form(signupFormEl, signupBtnEl);
  submit_form(signupFormEl, signupBtnEl, "/users/login");
}
if (window.location.href == "http://localhost:8080/users/login") {
  validate_form(loginFormEl, loginBtnEl);
  submit_form(loginFormEl, loginBtnEl, "/");
}
if (window.location.href == "http://localhost:8080/books/create") {
  validate_form(bookFormEl, bookBtnEl);
  submit_form(bookFormEl, bookBtnEl, "/books");
}

customElements.define("input-error", InputErrorElement);
