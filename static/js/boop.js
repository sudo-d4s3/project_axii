/* global Terminal */
/**
* Instance of Terminal
*/
(function (root, factory) {
  if (typeof module === 'object' && module.exports) {
    // Node. Does not work with strict CommonJS, but only CommonJS-like
    // environments that support module.exports, like Node.
    module.exports = factory(require('./terminal.js'), 'node')
  } else {
    // Browser globals (root is window)
    root.initTerminal = factory(root.Terminal, 'browser')
  }
}(this, function (Terminal, mode) {
  var commands = {}
  var state = {}

  commands.help = function() {
    var output = "<div>" +
      "<p><strong>Welcome Traveler!</strong> This is Project Axii, a place for me to build puzzels and for you to solve them!</p>" +
      "<br/>" +
      "<p>Your first goal is to find the password for the root account.</p>" +
      "<br/>" +
      "<p><em>Commands Available :</em><p>" +
      "<ul>" +
      "<li> echo: does what it says on the tin. ex 'echo &ltstuffz&gt' </li>" +
      "<li> su: changes the user. ex 'su &ltuser&gt &ltpassword&gt' </li>"
      "<li> </li>"
      "</ul>"
    return output
  }

  commands.echo = function (args) {
    args.shift()
    return args.join(' ')
  }

  commands.su = function (args) {
    if (args.length > 2) {
      if(args[1] === 'root'){
	if(args[2] === 'iHATEjs'){
          Terminal.user = args[1]
	} else {
	  return '<p><strong>Wrong Password!!!</strong></p>'
	}
      } else {
	return '<p><strong>User not found.</strong></p>'
      }
    }
    return ''
  }

  commands.exit = function (args) {
    Terminal.exit()
    console.log('[Process completed]')
  }

  function initTerminal() {
    console.log('Terminal access granted.')

    Terminal.init(document.body, {
      commands: commands,
      prompt: '\\u@\\H $ ',
      intro: '<p>Welcome to Project axii. Type \'help\' to get started.</p><p>&nbsp;</p>'
    })
  }

  // Return a value to define the module export.
  if (mode === 'node') {
    return {
      initTerminal: initTerminal,
      Terminal: Terminal
    }
  } else {
    return initTerminal
  }
}))
