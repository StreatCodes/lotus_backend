import { Nanocomponent } from "nanocomponent";
import { default as html } from "nanohtml";

export class Button extends Nanocomponent {
	color: String;
	constructor() {
		super();
	}

	createElement(color: String) {
		this.color = color;
		return html`
		<button style="background-color: ${color}">
			Click Me
		</button>
	`;
	}

	// Implement conditional rendering
	update(newColor: String) {
		return newColor !== this.color;
	}
}
