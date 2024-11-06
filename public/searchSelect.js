function searchSelect(button, dialog) {
    button.addEventListener("click", function() {
        dialog.showModal()
    })

    dialog.addEventListener("click" , function(event) {
        if (event.target === dialog) {
            dialog.close()
        }
    })
}
