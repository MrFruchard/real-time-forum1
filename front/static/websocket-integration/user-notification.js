export const notify = () => {
    let notifBubble = document.querySelector(".notif-bubble");

    // Si la bulle n'existe pas, la créer
    if (!notifBubble) {
        const btnNotify = document.querySelector(".notif");
        if (!btnNotify) return;

        notifBubble = document.createElement("div");
        notifBubble.classList.add("notif-bubble");
        notifBubble.style.position = "absolute";
        notifBubble.style.top = "1px";
        notifBubble.style.right = "2px";
        notifBubble.style.width = "10px";
        notifBubble.style.height = "10px";
        notifBubble.style.background = "red";
        notifBubble.style.borderRadius = "50%";
        notifBubble.style.display = "none"; // Cachée au départ

        btnNotify.appendChild(notifBubble);
        btnNotify.addEventListener("click", () => {
            notifBubble.remove()
        })
    }

    // Affiche la bulle rouge
    notifBubble.style.display = "block";
};