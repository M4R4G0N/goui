package components

import (
	"fmt"
)

type Command struct {
	Label  string
	Action string // Javascript para executar (ex: "window.location='/text'")
	Icon   string // Opcional
}

type CommandPaletteComponent struct {
	Commands []Command
	ID       string
}

func CommandPalette(cmds ...Command) *CommandPaletteComponent {
	return &CommandPaletteComponent{
		Commands: cmds,
		ID:       "goui-command-palette",
	}
}

func (c *CommandPaletteComponent) Render() string {
	// Gerar a lista de comandos para o JS
	cmdsJS := "["
	for i, cmd := range c.Commands {
		if i > 0 { cmdsJS += "," }
		cmdsJS += fmt.Sprintf("{label: '%s', action: function(){ %s }}", cmd.Label, cmd.Action)
	}
	cmdsJS += "]"

	return fmt.Sprintf(`
		<div id="%s" class="goui-command-palette-overlay" style="display:none">
			<div class="goui-command-palette-modal">
				<input type="text" id="%s-input" placeholder="Digite um comando..." class="goui-command-palette-input">
				<div id="%s-results" class="goui-command-palette-results"></div>
			</div>
		</div>

		<script>
			(function() {
				var overlay = document.getElementById('%s');
				var input = document.getElementById('%s-input');
				var results = document.getElementById('%s-results');
				var commands = %s;
				var selectedIndex = 0;

				function show() { overlay.style.display = 'flex'; input.focus(); render(); }
				function hide() { overlay.style.display = 'none'; }

				function render() {
					var filter = input.value.toLowerCase();
					var filtered = commands.filter(c => c.label.toLowerCase().includes(filter));
					results.innerHTML = '';
					filtered.forEach((c, i) => {
						var div = document.createElement('div');
						div.className = 'goui-command-palette-item' + (i === selectedIndex ? ' active' : '');
						div.innerText = c.label;
						div.onclick = function() { c.action(); hide(); };
						results.appendChild(div);
					});
				}

				document.addEventListener('keydown', function(e) {
					if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
						e.preventDefault();
						e.stopPropagation();
						show();
						return;
					}
					if (e.key === 'Escape') { hide(); return; }
					if (overlay.style.display === 'flex') {
						if (e.key === 'ArrowDown') { e.preventDefault(); selectedIndex++; render(); }
						if (e.key === 'ArrowUp') { e.preventDefault(); selectedIndex = Math.max(0, selectedIndex - 1); render(); }
						if (e.key === 'Enter') {
							var filter = input.value.toLowerCase();
							var filtered = commands.filter(c => c.label.toLowerCase().includes(filter));
							if (filtered[selectedIndex]) { filtered[selectedIndex].action(); hide(); }
						}
					}
				}, true);

				input.addEventListener('input', function() { selectedIndex = 0; render(); });
				overlay.onclick = function(e) { if (e.target === overlay) hide(); };
			})();
		</script>
	`, c.ID, c.ID, c.ID, c.ID, c.ID, c.ID, cmdsJS)
}
