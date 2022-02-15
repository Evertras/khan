# Components

Components are intended to be smaller, reusable pieces that might make up a
Screen.  They generally conform to `Init()`, `Update(tea.Msg)`, and `View()`
for simplicity, but `Update` returns a copy of itself rather than the more
generic `tea.Model`.  This allows for advanced controls/data to be passed around.

