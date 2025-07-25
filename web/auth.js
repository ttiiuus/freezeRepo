let isLogin = true

const emailField = document.getElementById('email-field')
const formTitle = document.getElementById('form-title')
const toggleLink = document.getElementById('toggle-link')
const toggleText = document.getElementById('toggle-text')
const submitBtn = document.getElementById('submit-btn')
const authForm = document.getElementById('auth-form')
const passwordInput = document.getElementById('password')
const showPasswordArea = document.getElementById('show-password-area')
const eyeIcon = document.getElementById('eye-icon')

// Переключение между Login и Register
toggleLink.addEventListener('click', () => {
	isLogin = !isLogin
	emailField.style.display = isLogin ? 'none' : 'block'
	formTitle.textContent = isLogin ? 'Login' : 'Register'
	submitBtn.textContent = isLogin ? 'Login' : 'Register'
	toggleText.textContent = isLogin
		? "Don't have an account?"
		: 'Already have an account?'
	toggleLink.textContent = isLogin ? 'Register' : 'Login'

	if (isLogin) {
		authForm.email.value = ''
	}
})

// При загрузке страницы проверяем JWT по куке
window.addEventListener('DOMContentLoaded', async () => {
	console.log('Проверка авторизации запускается...')
	try {
		const res = await fetch('http://192.168.209.1:8083/api/check', {
			credentials: 'include',
		})

		if (res.ok) {
			console.log('Авторизация успешна. Переход на upload...')
			window.location.href = 'http://192.168.209.1:8085'
		} else {
			console.log('Авторизация не пройдена.')
		}
	} catch (err) {
		console.error('Ошибка проверки авторизации:', err)
	}
})

// Submit формы (логин или регистрация)
authForm.addEventListener('submit', async e => {
	e.preventDefault()

	const username = authForm.username.value.trim()
	const password = authForm.password.value
	const email = authForm.email?.value.trim()

	const payload = isLogin
		? { username, password }
		: { username, password, email }

	const path = isLogin ? '/login' : '/register'

	try {
		const res = await fetch(`http://192.168.209.1:8083${path}`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload),
			credentials: 'include', // чтобы браузер сохранил Set-Cookie
		})

		if (res.ok) {
			window.location.href = 'http://192.168.209.1:8085'
		} else {
			alert('Ошибка: логин или регистрация не удалась')
		}
	} catch (err) {
		alert('Сервер недоступен')
	}
})

// Показать пароль при зажатии и скрыть при отпускании
showPasswordArea.addEventListener('mousedown', () => {
	passwordInput.type = 'text'
	eyeIcon.classList.replace('text-gray-400', 'text-blue-600')
})

showPasswordArea.addEventListener('mouseup', () => {
	passwordInput.type = 'password'
	eyeIcon.classList.replace('text-blue-600', 'text-gray-400')
})

showPasswordArea.addEventListener('mouseleave', () => {
	passwordInput.type = 'password'
	eyeIcon.classList.replace('text-blue-600', 'text-gray-400')
})
