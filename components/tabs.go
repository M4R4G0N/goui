package components

import (
	"fmt"
	"strings"
)

type TabItem struct {
	Text     string
	Href     string
	Active   bool
	TargetID string // ID da div que será mostrada ao clicar
}

type TabsComponent struct {
	ID    string
	Items []TabItem
}

func (t *TabsComponent) Render() string {
	if t.ID == "" {
		t.ID = AutoID()
	}

	tabsHTML := ""
	targets := []string{}

	for _, item := range t.Items {
		activeClass := ""
		if item.Active {
			activeClass = "active"
		}

		// Se não houver Href, usamos '#' para não recarregar
		href := item.Href
		if href == "" {
			href = "javascript:void(0)"
		}

		tabsHTML += fmt.Sprintf(
			`<a href="%s" class="goui-tab %s" onclick="switchTab(this, '%s', '%s')">%s</a>`,
			href, activeClass, t.ID, item.TargetID, item.Text,
		)

		if item.TargetID != "" {
			targets = append(targets, item.TargetID)
		}
	}

	// Script para gerenciar a troca de abas e visibilidade
	script := fmt.Sprintf(`
		<script>
			if (typeof window.switchTab !== 'function') {
				window.switchTab = function(el, groupId, targetId) {
					// 1. Remove active de todas as abas do mesmo grupo
					var container = el.closest('.goui-tabs');
					container.querySelectorAll('.goui-tab').forEach(t => t.classList.remove('active'));
					
					// 2. Adiciona active na clicada
					el.classList.add('active');
					
					// 3. Gerencia visibilidade dos alvos (se existirem)
					if (targetId) {
						// Aqui a gente assume que todos os possíveis alvos devem ser escondidos
						// Vou buscar os alvos baseados no que foi passado no componente
						var allTargets = %s;
						allTargets.forEach(id => {
							var target = document.getElementById(id);
							if (target) target.style.display = 'none';
						});
						
						var target = document.getElementById(targetId);
						if (target) target.style.display = 'block';
					}
				};
			}
		</script>
	`, fmt.Sprintf("['%s']", strings.Join(targets, "','")))

	return fmt.Sprintf(`<div class="goui-tabs" id="%s">%s</div>%s`, t.ID, tabsHTML, script)
}

func Tabs(items ...TabItem) *TabsComponent {
	return &TabsComponent{Items: items}
}
