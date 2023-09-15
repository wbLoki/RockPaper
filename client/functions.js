class GameState {
    constructor() {
        this.value = 'offline';
    }
    onChange(newState) {
        if (newState == "gameon") {
            document.getElementsByTagName("body")[0].className = "gameon"
            var phaseElement = document.getElementsByClassName("ph1")[0]
            if (phaseElement) {
                phaseElement.className = "hide"
            }
        }
    }

    set(newValue) {
        if (this.value !== newValue) {
            this.value = newValue;
            if (typeof this.onChange === 'function') {
                this.onChange(newValue);
            }
        }
    }

    get() {
        return this.value;
    }
}

// Function to set a cookie
function setCookie(name, value, daysToExpire) {
    const expirationDate = new Date();
    expirationDate.setDate(expirationDate.getDate() + daysToExpire);

    const cookieValue = `${name}=${encodeURIComponent(value)}; expires=${expirationDate.toUTCString()}; path=/`;

    document.cookie = cookieValue;
}