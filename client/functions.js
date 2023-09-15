class GameState {
    constructor() {
        this.value = 'offline';
    }
    onChange(newState) {
        if (newState == "gameon") {
            document.getElementsByTagName("body")[0].className = "gameon"
            document.getElementsByClassName("ph1")[0].className = "hide"
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
