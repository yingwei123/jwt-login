class SignUp{
    emailInput = "";
    passwordInput = "";
    confirmPasswordInput = "";

    constructor(){
        this.initHandler();
    }

    initHandler(){
        document.getElementById("signup-button").onclick = e => this.CreateNewUser();
        this.emailInput = document.getElementById("email");
        this.passwordInput = document.getElementById("password");
        this.confirmPasswordInput = document.getElementById("cPassword");
    }

    ValidUserParameters(){
        const currentEmail = this.emailInput.value;
        const currentPassword = this.passwordInput.value;
        const confirmPassword = this.confirmPasswordInput.value;

        if(this.ValidateEmail(currentEmail) === false){
            alert("Incorrect email format! Please try again");
            this.emailInput.value = "";
            return;
        }

        if(currentPassword.length < 8){
            alert("Make sure your password length is above 8!");
            this.passwordInput.value = "";
            return;
        }

        if(currentPassword != confirmPassword){
            alert("confirm password do not match! Please try again");
            this.confirmPasswordInput.value = "";
            return;
        }

        

        return [currentEmail, currentPassword];
    }

    ValidateEmail(mail) {
        if (/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(mail)){
            return (true);
        }
        return (false);
    }

    CreateNewUser(){
        const signupInfo = this.ValidUserParameters();
        if(typeof(signupInfo[0])==="undefined" || typeof(signupInfo[1]) ==="undefined"){
            return
        }

        const data = { email: signupInfo[0],
                       password:signupInfo[1],
        };

        console.log(data)

        fetch('/new-user', {
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

            window.location.replace("/login");
        })
        .catch((error) => {
        console.error('Error:', error);
        });
    }
}