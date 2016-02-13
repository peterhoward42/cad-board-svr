/* A thing to which you can send mouse move events to to be processed.
   It avoids needless or excessive processing by ignoring very small
   movements.

   Beware of challenges in trying to have the constructor register the html5 event callback itself.
   This doesn't work because the browser's runtime call cannot provide any object context to
   the call.
 */
var CadboardMouseTracker = function(sourceOfMouseMovements) {
    // Constants
    this.prevX = 1e6; // out of band sentinel value
    this.prevY = 1e6;
    this.negligibleMovementSize = 5;

    // Cache geometry of html element doing the sending to avoid recalculating it for
    // each event received.
    this.rect = sourceOfMouseMovements.getBoundingClientRect();
    this.width = this.rect.right - this.rect.left;
    this.height = this.rect.bottom - this.rect.top;
}

CadboardMouseTracker.prototype.handleMoveEvent = function(mouseEvent) {
    var mouseX = mouseEvent.clientX - this.rect.left;
    var mouseY = mouseEvent.clientY - this.rect.top;
    // Optimise out moves that are too small (or zero) while still in integer coordinates
    if (this.movementTooSmall(mouseX, mouseY))
        return;
    this.prevX = mouseX;
    this.prevY = mouseY;
    // Normalise the coords
    mouseX = mouseX / this.width;
    mouseY = mouseY / this.height;
};

CadboardMouseTracker.prototype.movementTooSmall = function(x, y) {
    if (Math.abs(x - this.prevX) > this.negligibleMovementSize)
        return false;
    if (Math.abs(y - this.prevY) > this.negligibleMovementSize)
        return false;
    return true;
};