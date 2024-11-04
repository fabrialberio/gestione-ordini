function searchSelect(input, dialog, searchInput, options, select) {
    const items = options.children
    
    function onSearchInput() {
        const value = searchInput.value.toLowerCase()

        for (let i = 0; i < items.length; i++) {
            const item = items[i]
            const text = item.innerHTML.toLowerCase()

            if (text.includes(value)) {
                item.style.display = "block"
            } else {
                item.style.display = "none"
            }
        }
    }

    function openDialog() {
        dialog.style.minWidth = input.offsetWidth + "px"
        dialog.show()
    }

    function onItemClick(event) {
        const item = event.target
        select.value = item.value
        input.innerHTML = item.innerHTML
        dialog.close()
    }

    searchInput.addEventListener("input", onSearchInput)
    input.addEventListener("click", openDialog)

    for (let i = 0; i < items.length; i++) {
        const item = items[i]

        item.addEventListener("click", onItemClick)
    }

    input.innerHTML = select.options[select.selectedIndex].innerHTML
}
