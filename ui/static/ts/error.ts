export class InputErrorElement extends HTMLElement {
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
    } else {
      this.setAttribute("value", newValue);
    }
  }

  get for() {
    return this.getAttribute("for");
  }

  connectedCallback() {
    this.textContent = this.getAttribute("value") || "";
  }

  attributeChangedCallback(name: any, _oldValue: any, newValue: any) {
    if (name == "value") {
      this.textContent = newValue;
    }
  }
}
