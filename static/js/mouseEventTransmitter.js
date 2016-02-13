/*
This object will track mouse movement events from an given html5 element, and relay them
with a POST message to a an http server. (asynchronously).

The constructor takes care of registering itself as an event handler for the html element you
specify.

It takes some steps to optimise the communications. Firstly it ignores small mouse movements.
Secondly it throttles its sending of messages to match the rate at which the server can receive
and reply to them. Any mouse movements it receives while waiting for the server to be ready are
discarded, so the server receives the most recently available position always.
*/

// Factory method for the object.
function mouseEventTransmitter(mouseEventSource) {
    t = {};
    // Constants
    t.rect = mouseEventSource.getBoundingClientRect();
    t.left = t.rect.left;
    t.top = t.rect.top;
    t.width = t.rect.right - t.rect.left;
    t.height = t.rect.bottom - t.rect.top;
    t.serverUrl = "/mouse";

    t.negligibleMovementSize = 5;

    // State
    t.prevX = 1e6; // An out of band sentinel value.
    t.prevY = 1e6;
    t.waitingForServerReply = false;

    // Functions

    // Handler to receive mouse move events from the browser
    t.handleMoveEvent= function(event) {
        mouseX = event.clientX - t.left;
        mouseY = event.clientY - t.top;
        // Optimise out moves that are too small (or zero) while still in integer coordinates
        if (t.movementTooSmall(mouseX, mouseY))
            return;
        t.prevX = mouseX;
        t.prevY = mouseY;

        // Throttle the message sending rate to match the server's replies.
        if (t.waitingForServerReply)
            return;

        // Send the new position to the server
        normalisedX = mouseX / t.width;
        normalisedY = mouseY / t.height;
        t.postPositionToServerAsync(normalisedX, normalisedY);
    };

    // Utility function to assess if mouse movement is very small or zero
    t.movementTooSmall = function(x, y) {
          if (Math.abs(x - t.prevX) > t.negligibleMovementSize)
              return false;
          if (Math.abs(y - t.prevY) > t.negligibleMovementSize)
              return false;
          return true;
     };

    // Utility to post an x / y position to the server.
    t.postPositionToServerAsync = function(x, y) {
        payload = {X: normalisedX, Y: normalisedY};
        callbackFn = t.handlerForServerReply;
        $.post(t.serverUrl, payload, callbackFn);
    }

    // Handler to accept replies from the server.
    t.handlerForServerReply = function(reply) {
        t.waitingForServerReply = false;
    }

    // Ask the mouse move event emitter to send the events here.
    mouseEventSource.onmousemove = t.handleMoveEvent;
    return t;
}