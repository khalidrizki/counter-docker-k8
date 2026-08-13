console.log("SCRIPT TERLOAD"); 
const button = document.getElementById("increment-button");
const countElement = document.getElementById("count");

async function loadCount() {
    const response = await fetch("/api/count");

    const data = await response.json();

    countElement.innerText = data.count;
}

loadCount();

button.addEventListener("click", async function() {
    const response = await fetch("/api/count", {
        method:"POST"
    });

    const data = await response.json();
    console.log(data)
    countElement.innerText = data.count;
});