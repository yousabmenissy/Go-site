import { InputErrorElement } from "./error";

export function validate_form(form: HTMLFormElement, btn: HTMLButtonElement) {
  const elements = Array.from(form.elements) as (
    | HTMLInputElement
    | HTMLTextAreaElement
  )[];

  const passwordEl = document.querySelector(
    "input[name=password]"
  ) as HTMLInputElement;

  const err = document.querySelectorAll(
    "input-error"
  ) as NodeListOf<InputErrorElement>;

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
          } else {
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
    } else {
      btn.disabled = true;
    }
  });
}

export function submit_form(
  form: HTMLFormElement,
  btn: HTMLButtonElement,
  addr: string
) {
  const elements = Array.from(form.elements) as (
    | HTMLInputElement
    | HTMLTextAreaElement
  )[];

  const err = document.querySelectorAll(
    "input-error"
  ) as NodeListOf<InputErrorElement>;

  form.addEventListener("submit", (event) => {
    event.preventDefault();

    const formData: FormData = new FormData(form, btn);
    const data: { [key: string]: FormDataEntryValue } = {};
    formData.forEach((value, key) => {
      data[key] = value;
    });

    const csrfToken = document.querySelector(
      "input[name=csrf_token]"
    ) as HTMLInputElement;
    fetch(btn.getAttribute("url")!, {
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
        for (const [key, val] of Object.entries<string>(data)) {
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
