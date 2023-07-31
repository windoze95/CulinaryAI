document.querySelectorAll('.toggle-symbols').forEach(function (icon) {
    icon.addEventListener('click', function () {
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

document.getElementById('button1').addEventListener('click', function () {
});

document.getElementById('button2').addEventListener('click', function () {
});

// Markdown processing function
function processMarkdown(text) {
    var html = marked(text);
    document.getElementById('markdown-display').innerHTML = html;
}

$(document).ready(function () {
    $(".dropdown-trigger").dropdown();
});

document.addEventListener('DOMContentLoaded', function () {
    var elems = document.querySelectorAll('.sidenav');
    var html = document.documentElement;
    var body = document.body;
    var instances = M.Sidenav.init(elems, {
        edge: 'right',
        draggable: true,
        preventScrolling: true,
        onOpenStart: function () {
            // document.getElementById('menu-icon').textContent = 'close';
            body.classList.add('overlay-active');
        },
        onCloseStart: function () {
            body.classList.remove('overlay-active');
        },
        onCloseEnd: function () {
            // document.getElementById('menu-icon').textContent = 'menu';  // Change icon back to 'menu' when sidenav is closed
        }
    });
    var elemsTabs = document.querySelectorAll('.tabs');
    var instancesTabs = M.Tabs.init(elemsTabs);
});

// Handle the Settings form submission
document.getElementById('settingsForm').addEventListener('submit', function (e) {
    e.preventDefault();

    var apiKey = document.getElementById('apiKeyInput').value;

    // Create a new FormData object
    var formData = new FormData();
    formData.append('OpenAIKey', apiKey);

    // Send a POST request to the server
    fetch('/settings', {
        method: 'POST',
        body: formData
    }).then(response => {
        // Check if the request was successful
        if (response.ok) {
            console.log('OpenAI API Key was saved successfully');
        } else {
            console.error('Failed to save the OpenAI API Key');
        }
    }).catch(error => {
        console.error('An error occurred:', error);
    });

    // Clear the input field and hide the modal
    document.getElementById('apiKeyInput').value = '';
    modal.style.display = "none";
});

window.scrollTo(0, 1);
