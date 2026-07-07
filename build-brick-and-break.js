let brickCount = 0;
let timer = null;

const getTower = () => document.getElementById("tower") || document.body;

export function build(amount) {
    clearInterval(timer);

    let built = 0;

    timer = setInterval(() => {
        if (built >= amount) {
            clearInterval(timer);
            return;
        }

        brickCount++;

        const brick = document.createElement("div");
        brick.id = `brick-${brickCount}`;

        // Every second brick in a row of three belongs to the middle column
        if (brickCount % 3 === 2) {
            brick.dataset.foundation = "true";
        }

        getTower().appendChild(brick);

        built++;
    }, 100);
}

export function repair(...ids) {
    ids.forEach((id) => {
        const brick = document.getElementById(id);
        if (!brick) return;

        if (brick.dataset.foundation === "true") {
            brick.dataset.repaired = "in progress";
        } else {
            brick.dataset.repaired = "true";
        }
    });
}

export function destroy() {
    const tower = getTower();
    const bricks = tower.querySelectorAll('div[id^="brick-"]');

    if (bricks.length === 0) return;

    bricks[bricks.length - 1].remove();
    brickCount--;
}
