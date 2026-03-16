(function () {
    var links = document.querySelectorAll('.site-nav a[href^="#"]');
    var sections = ['info', 'units', 'contact', 'schedule']
        .map(function (id) {
            return document.getElementById(id);
        })
        .filter(Boolean);

    function setActiveLink() {
        var current = sections[0] ? sections[0].id : '';
        var top = window.scrollY + 140;

        sections.forEach(function (section) {
            if (top >= section.offsetTop) {
                current = section.id;
            }
        });

        links.forEach(function (link) {
            var isActive = link.getAttribute('href') === '#' + current;
            link.classList.toggle('active', isActive);
        });
    }

    window.addEventListener('scroll', setActiveLink, { passive: true });
    window.addEventListener('load', setActiveLink);
})();
