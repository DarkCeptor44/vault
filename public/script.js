const d = (x) => document.querySelector(x);
const dA = (x) => document.querySelectorAll(x);

if (typeof (alertify) !== "undefined") {
    window.alert = alertify.log;
}

/**
 * 
 * @param {Element[]} fields 
 * @returns {boolean}
 */
function checkFields(fields) {
    let ret = false;
    for (let field of fields) {
        if (field.classList.contains("error")) field.classList.remove("error");
        ret = field.type === 'number' ? field.value === 0 : field.value.trim() === "";
        if (ret) {
            field.classList.add("error");
        }
    }
    return ret;
}

Object.defineProperty(String.prototype, "capitalize", {
    value: function () {
        return this.charAt(0).toUpperCase() + this.slice(1);
    },
    enumerable: false
});