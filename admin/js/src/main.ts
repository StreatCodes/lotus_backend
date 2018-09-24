
import * as API from './api';

const loginForm = document.getElementById('loginForm');
const loginErrorBox = document.getElementById('loginErrorBox');
const loginSubmitButton = document.getElementById('loginSubmitButton');

loginForm.addEventListener('submit', async (e) => {
		e.preventDefault();

		let email: String = loginForm.children['email'].value;
		let password: String = loginForm.children['password'].value;

		loginErrorBox.classList.add('flat');
		loginErrorBox.innerHTML = '';
		loginSubmitButton.classList.add('loading');


		try {
			let result = await API.login(email, password);
			loginSubmitButton.classList.remove('loading');

			localStorage.setItem("session", JSON.stringify(result));
			window.location.assign('/admin/');
		} catch(e) {
			loginErrorBox.classList.remove('flat');
			loginErrorBox.innerText = e.message;
			loginSubmitButton.classList.remove('loading');
		}
	});
