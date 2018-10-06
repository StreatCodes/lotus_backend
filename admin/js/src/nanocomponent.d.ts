declare module 'nanocomponent' {
    export class Nanocomponent {
        constructor(name?: String);

        render(args: any): HTMLElement;
        rerender(): void;
        _handleRender(args: any): HTMLElement;
        _createProxy(): HTMLElement;
        _reset(): void;
        _brandNode(node: HTMLElement): HTMLElement;
        _ensureID(node: HTMLElement): HTMLElement;
        _handleLoad(el: HTMLElement): void;
        _handleUnload(el: HTMLElement): void;
        createElement(args: any): HTMLElement;
        update(args: any): Boolean;
    }
}

// Nanocomponent.prototype.render = function () {
//   var renderTiming = nanotiming(this._name + '.render')
//   var self = this
//   var args = new Array(arguments.length)
//   var el
//   for (var i = 0; i < arguments.length; i++) args[i] = arguments[i]
//   if (!this._hasWindow) {
//     var createTiming = nanotiming(this._name + '.create')
//     el = this.createElement.apply(this, args)
//     createTiming()
//     renderTiming()
//     return el
//   } else if (this.element) {
//     el = this.element // retain reference, as the ID might change on render
//     var updateTiming = nanotiming(this._name + '.update')
//     var shouldUpdate = this._rerender || this.update.apply(this, args)
//     updateTiming()
//     if (this._rerender) this._rerender = false
//     if (shouldUpdate) {
//       var desiredHtml = this._handleRender(args)
//       var morphTiming = nanotiming(this._name + '.morph')
//       morph(el, desiredHtml)
//       morphTiming()
//       if (this.afterupdate) this.afterupdate(el)
//     }
//     if (!this._proxy) { this._proxy = this._createProxy() }
//     renderTiming()
//     return this._proxy
//   } else {
//     this._reset()
//     el = this._handleRender(args)
//     if (this.beforerender) this.beforerender(el)
//     if (this.load || this.unload || this.afterreorder) {
//       onload(el, self._handleLoad, self._handleUnload, self._ncID)
//       this._olID = el.dataset[OL_KEY_ID]
//     }
//     renderTiming()
//     return el
//   }
// }

// Nanocomponent.prototype.rerender = function () {
//   assert(this.element, 'nanocomponent: cant rerender on an unmounted dom node')
//   this._rerender = true
//   this.render.apply(this, this._arguments)
// }

// Nanocomponent.prototype._handleRender = function (args) {
//   var createElementTiming = nanotiming(this._name + '.createElement')
//   var el = this.createElement.apply(this, args)
//   createElementTiming()
//   if (!this._rootNodeName) this._rootNodeName = el.nodeName
//   assert(el instanceof window.Element, 'nanocomponent: createElement should return a DOM node')
//   assert.equal(this._rootNodeName, el.nodeName, 'nanocomponent: root node types cannot differ between re-renders')
//   this._arguments = args
//   return this._brandNode(this._ensureID(el))
// }

// Nanocomponent.prototype._createProxy = function () {
//   var proxy = document.createElement(this._rootNodeName)
//   var self = this
//   this._brandNode(proxy)
//   proxy.id = this._id
//   proxy.setAttribute('data-proxy', '')
//   proxy.isSameNode = function (el) {
//     return (el && el.dataset.nanocomponent === self._ncID)
//   }
//   return proxy
// }

// Nanocomponent.prototype._reset = function () {
//   this._ncID = makeID()
//   this._olID = null
//   this._id = null
//   this._proxy = null
//   this._rootNodeName = null
// }

// Nanocomponent.prototype._brandNode = function (node) {
//   node.setAttribute('data-nanocomponent', this._ncID)
//   if (this._olID) node.setAttribute(OL_ATTR_ID, this._olID)
//   return node
// }

// Nanocomponent.prototype._ensureID = function (node) {
//   if (node.id) this._id = node.id
//   else node.id = this._id = this._ncID
//   // Update proxy node ID if it changed
//   if (this._proxy && this._proxy.id !== this._id) this._proxy.id = this._id
//   return node
// }

// Nanocomponent.prototype._handleLoad = function (el) {
//   if (this._loaded) {
//     if (this.afterreorder) this.afterreorder(el)
//     return // Debounce child-reorders
//   }
//   this._loaded = true
//   if (this.load) this.load(el)
// }

// Nanocomponent.prototype._handleUnload = function (el) {
//   if (this.element) return // Debounce child-reorders
//   this._loaded = false
//   if (this.unload) this.unload(el)
// }

// Nanocomponent.prototype.createElement = function () {
//   throw new Error('nanocomponent: createElement should be implemented!')
// }

// Nanocomponent.prototype.update = function () {
//   throw new Error('nanocomponent: update should be implemented!')
// }
