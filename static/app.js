document.querySelectorAll('.toggle-symbols').forEach(function(icon) {
    icon.addEventListener('click', function() {
        switch (this.textContent) {
            case 'favorite':
                this.textContent = 'favorite_border';
                break;
            case 'favorite_border':
                this.textContent = 'favorite';
                break;
            case 'delete':
                this.textContent = 'auto_delete';
                break;
            case 'auto_delete':
                this.textContent = 'delete';
                break;
            default:
                break;
        }
    });
});

// Markdown processing function
function processMarkdown(text) {
    var html = marked(text);
    document.getElementById('markdown-display').innerHTML = html;
}

$(document).ready(function() {
    $(".dropdown-trigger").dropdown();
});

document.addEventListener('DOMContentLoaded', function() {
    var elems = document.querySelectorAll('.sidenav');
    var html = document.documentElement;
    var body = document.body;
    var instances = M.Sidenav.init(elems, {
        edge: 'right',
        draggable: true,
        preventScrolling: true,
        onOpenStart: function() {
            // document.getElementById('menu-icon').textContent = 'close';
            body.classList.add('overlay-active');
        },
        onCloseStart: function() {
            body.classList.remove('overlay-active');
        },
        onCloseEnd: function() {
            // document.getElementById('menu-icon').textContent = 'menu';  // Change icon back to 'menu' when sidenav is closed
        }
    });
    var elemsTabs = document.querySelectorAll('.tabs');
    var instancesTabs = M.Tabs.init(elemsTabs);
});

document.querySelector("#generate-recipe-button").addEventListener("click", function(e) {
    e.preventDefault(); // Prevent the default form submission

    // Make the fetch call
    fetch("/recipes", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                prompt: document.querySelector("#user-prompt-input").value,
            }),
        })
        .then(async response => {
            const contentType = response.headers.get("content-type");
            if (contentType && contentType.includes("application/json")) {
                const data = await response.json();
                if (!response.ok) {
                    console.log(data.error);
                    M.toast({ html: data.error || "An error occurred" });
                    throw new Error(data.error || response.statusText);
                }
                return data;
            } else {
                console.log(response.text());
                M.toast({ html: "An error occurred" });
                throw new Error("Invalid content type");
            }
        })
        .then(data => {
            // Insert the recipe (markdown) into an element on the same page
            document.querySelector("#markdown-display").innerHTML = data.recipe;
        })
        .catch(error => {
            // Handle any other errors
            if (error.message === "") {
                M.toast({ html: "An unknown error occurred" });
            } else {
                M.toast({ html: error.message });
            }
        });
});

// document.querySelector("#generate-recipe-button").addEventListener("click", function(e) {
//     e.preventDefault(); // Prevent the default form submission

//     // Make the fetch call
//     fetch("/recipes", {
//             method: "POST",
//             headers: {
//                 "Content-Type": "application/json"
//             },
//             body: JSON.stringify({
//                 userPrompt: document.querySelector("#user-prompt-input").value,
//             }),
//         })
//         .then(response => {
//             if (!response.ok) {
//                 // const data = await response.json();
//                 console.log(data.error)
//                 throw new Error(response.statusText);
//             }
//             return response.json();
//             //     if (!response.ok) {
//             //     // Check if the response has a JSON content type
//             //     const contentType = response.headers.get("content-type");
//             //     if (contentType && contentType.includes("application/json")) {
//             // Parse the JSON response
//             // const data = await response.json();
//             //         // Display the error message using Materialize toast
//             //         M.toast({ html: data.error || "An error occurred" });
//             //         throw new Error("Server error");
//             //     } else {
//             //         // If not JSON, just throw an error
//             //         throw new Error("Server error");
//             //     }

//             // Handle error response from server
//             // response.json().then((data) => {
//             // Display the error message using Materialize toast
//             // M.toast({ html: data.error || "An error occurred" });
//         })
//         // throw new Error("Server error");

//     // }
//     // return response.json();
//     // })
//     .then(data => {
//             // Insert the recipe (markdown) into an element on the same page
//             document.querySelector("#markdown-display").innerHTML = data.recipe;
//         })
//         .catch(error => {
//             // Handle any other errors
//             if (error.message === "") {
//                 M.toast({ html: "an unknown error occurred" })
//             } else {
//                 M.toast({ html: error.message })
//             }
//         });
// });

document.getElementById('openSettings').addEventListener('click', function(e) {
    e.preventDefault(); // Prevent the default link behavior

    fetch('/settings')
        .then(response => response.text())
        .then(data => {
            // Update the content of the modal
            document.getElementById('settingsContent').innerHTML = data;

            var modalInstance = M.Modal.getInstance(document.getElementById('settingsModal'));
            if (!modalInstance) {
                modalInstance = M.Modal.init(document.getElementById('settingsModal'));
            }
            modalInstance.open();
        })
        .catch(error => M.toast({ html: error.message }));
});

document.getElementById('saveSettings').addEventListener('click', function() {
    // Retrieve the API key from the form
    var apikey = document.getElementById('apikey').value;

    // Send the PUT request to the server
    fetch('/users/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ apikey: apikey })
        })
        .then(response => response.json())
        .then(data => {
            // Assuming the response has a 'message' property
            if (data.message) {
                M.toast({ html: data.message });
            }
        })
        .catch(error => M.toast({ html: error.message }));

    // Close the modal
    var modalInstance = M.Modal.getInstance(document.getElementById('settingsModal'));
    modalInstance.close();
});

// Handle the Settings form submission
// document.getElementById('settingsForm').addEventListener('submit', function (e) {
//     e.preventDefault();

//     var apiKey = document.getElementById('apiKeyInput').value;

//     // Create a new FormData object
//     var formData = new FormData();
//     formData.append('OpenAIKey', apiKey);

//     // Send a POST request to the server
//     fetch('/settings', {
//         method: 'POST',
//         body: formData
//     }).then(response => {
//         // Check if the request was successful
//         if (response.ok) {
//             console.log('OpenAI API Key was saved successfully');
//         } else {
//             console.error('Failed to save the OpenAI API Key');
//         }
//     }).catch(error => {
//         console.error('An error occurred:', error);
//     });

//     // Clear the input field and hide the modal
//     document.getElementById('apiKeyInput').value = '';
//     modal.style.display = "none";
// });

// window.scrollTo(0, 1);