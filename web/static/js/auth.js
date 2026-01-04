document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('.auth-form');

    form.addEventListener('submit', (e) => {
        const password = form.querySelector('input[name="password"]').value;
        
        if (password.length < 6) {
            e.preventDefault();
            alert('Пароль должен быть не менее 6 символов!');
        }
    });
});