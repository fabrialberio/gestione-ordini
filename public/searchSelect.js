function searchSelect(input, dialog, searchInput, options, select) {
    const items = options.children

    function onSearchInput() {
        const value = searchInput.value.toLowerCase()

        for (let i = 0; i < items.length; i++) {
            const item = items[i]
            const text = item.innerHTML.toLowerCase()

            if (text.includes(value)) {
                item.style.display = ""
            } else {
                item.style.display = "none"
            }
        }
    }
    searchInput.addEventListener("input", onSearchInput)
    
    input.addEventListener("click", function() {
        dialog.showModal()
    })

    dialog.addEventListener("click" , function(event) {
        if (event.target === dialog) {
            dialog.close()
        }
    })

    for (let i = 0; i < items.length; i++) {
        const item = items[i]

        item.addEventListener("click", function(event) {
            const item = event.target
            select.value = item.value
            input.innerHTML = item.innerHTML
            dialog.close()
        })
    }

    input.innerHTML = select.options[select.selectedIndex].innerHTML
}
