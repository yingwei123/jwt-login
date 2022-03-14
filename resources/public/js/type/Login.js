class Login{
    emailInput ="";
    passwordInput = "";

    constructor(){
        this.initHandler();
    }

    initHandler(){
        document.getElementById("login").onclick = e => this.handleLoginClick()
        this.emailInput = document.getElementById("login-email");
        this.passwordInput = document.getElementById("login-password");
    }

    handleLoginClick(){
        const data = { email: this.emailInput.value,
            password: this.passwordInput.value,
        };

        console.log(data)

        fetch('/auth', {
            method: 'POST', // or 'PUT'
            headers: {
            'Content-Type': 'application/json',
        },
            body: JSON.stringify(data),
        })
        .then(response => response.json())
        .then(data => {
            if(data.Status != 200){
                alert(data.Message)
                return
            }

            console.log(data)

            window.location.replace("/default");
            return
        })
        .catch((error) => {
        console.error('Error:', error);
        });
    }
}