function validate_form(form, btn) {
    const elements = Array.from(form.elements);
    const passwordEl = document.querySelector("input[name=password]");
    const err = document.querySelectorAll("input-error");
    elements.forEach((input) => {
        input.addEventListener("blur", () => {
            err.forEach((err) => {
                if (err.for == input.name) {
                    err.value = input.validationMessage;
                }
            });
            input.classList.add("was-validated");
        });
        input.addEventListener("input", () => {
            input.removeAttribute("class");
            if (passwordEl) {
                if (input.name == "password-confirm") {
                    if (input.value !== passwordEl.value) {
                        input.setCustomValidity("password mismatch");
                    }
                    else {
                        input.setCustomValidity("");
                    }
                }
            }
            err.forEach((err) => {
                if (err.for == input.name) {
                    err.value = "";
                }
            });
        });
    });
    form.addEventListener("input", () => {
        if (form.checkValidity()) {
            btn.disabled = false;
        }
        else {
            btn.disabled = true;
        }
    });
}
function submit_form(form, btn, addr) {
    const elements = Array.from(form.elements);
    const err = document.querySelectorAll("input-error");
    form.addEventListener("submit", (event) => {
        event.preventDefault();
        const formData = new FormData(form, btn);
        const data = {};
        formData.forEach((value, key) => {
            data[key] = value;
        });
        const csrfToken = document.querySelector("input[name=csrf_token]");
        fetch(btn.getAttribute("url"), {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-CSRF-Token": csrfToken.value,
            },
            body: JSON.stringify(data),
        })
            .then((response) => {
            if (response.ok) {
                window.location.href = addr;
            }
            return response.json();
        })
            .then((data) => {
            console.log("Response Data:", data);
            for (const [key, val] of Object.entries(data)) {
                elements.forEach((input) => {
                    err.forEach((err) => {
                        if (input.name == key && err.for == key) {
                            input.setCustomValidity(val);
                            err.value = val;
                            input.classList.add("was-validated");
                        }
                    });
                });
            }
        })
            .catch((err) => {
            console.error(err);
        });
    });
}

class InputErrorElement extends HTMLElement {
    static observedAttributes = ["value"];
    constructor() {
        super();
    }
    get value() {
        return this.getAttribute("value");
    }
    set value(newValue) {
        if (newValue === null || newValue == "") {
            this.removeAttribute("value");
        }
        else {
            this.setAttribute("value", newValue);
        }
    }
    get for() {
        return this.getAttribute("for");
    }
    connectedCallback() {
        this.textContent = this.getAttribute("value") || "";
    }
    attributeChangedCallback(name, _oldValue, newValue) {
        if (name == "value") {
            this.textContent = newValue;
        }
    }
}

const signupFormEl = document.querySelector("#signup form");
const signupBtnEl = document.querySelector("#signup button[type=submit]");
const loginFormEl = document.querySelector("#login form");
const loginBtnEl = document.querySelector("#login button[type=submit]");
const bookFormEl = document.querySelector("#book-create form");
const bookBtnEl = document.querySelector("#book-create button[type=submit]");
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
