export function login(email: String, password: String): Promise<Object> {
	return fetch("/api/authorize", {
		method: "POST",
		body: JSON.stringify({
			email: email,
			password: password
		})
	}).then(async res => {
		if(res.status === 200) {
			return res.text();
		} else {
			throw new Error(await res.text());
		}
	});
}