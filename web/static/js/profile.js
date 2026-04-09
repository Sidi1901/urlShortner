  
function showTab(tab) {
    const contents = document.querySelectorAll('.tab-content');
    const buttons = document.querySelectorAll('.tab-btn');

    contents.forEach(c => c.classList.remove('active'));
    buttons.forEach(b => b.classList.remove('active'));

    document.getElementById(tab).classList.add('active');
    event.target.classList.add('active');
}
