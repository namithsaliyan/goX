document.getElementById('contactForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const formData = new FormData(this);
    const data = {
        name: formData.get('name'),
        email: formData.get('email'),
        message: formData.get('message')
    };

    fetch('https://38c4a8e5-f2cc-498c-b22b-237261dbdd0a-00-340wwgf6e0cpa.sisko.replit.dev/submit', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })
    .then(response => response.json())
    .then(result => {
        alert(result.message);
    })
    .catch(error => {
        console.error('Error:', error);
    });
});
document.addEventListener('DOMContentLoaded', function() {
    const gallery = document.querySelector('.gallery');
    const items = document.querySelectorAll('.gallery-item');

    items.forEach(item => {
        item.addEventListener('click', () => {
            gallery.scrollTo({
                left: item.offsetLeft - (gallery.clientWidth / 2) + (item.clientWidth / 2),
                behavior: 'smooth'
            });
        });
    });
});

let lastScrollTop = 0;
const header = document.querySelector('header');

window.addEventListener('scroll', () => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;

    if (scrollTop > lastScrollTop) {
        // Scrolling down
        header.classList.add('hidden');
    } else {
        // Scrolling up
        header.classList.remove('hidden');
    }
    lastScrollTop = scrollTop <= 0 ? 0 : scrollTop; // For Mobile or negative scrolling
});


document.addEventListener('DOMContentLoaded', function() {
    const header = document.querySelector('header');
    const sections = document.querySelectorAll('section');
    const navLinks = document.querySelectorAll('nav ul li a');

    function changeHeaderBackground() {
        let index = sections.length;

        while (--index && window.scrollY + 50 < sections[index].offsetTop) {}

        navLinks.forEach(link => link.classList.remove('active'));
        if (index >= 0) {
            navLinks[index].classList.add('active');
            const sectionBackground = sections[index].style.backgroundImage || '';
            header.style.background = sectionBackground;
        }
    }

    window.addEventListener('scroll', changeHeaderBackground);
    changeHeaderBackground(); // Call on page load
});
