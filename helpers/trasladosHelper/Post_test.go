package trasladoshelper

import (
	"strings"
	"testing"

	"github.com/udistrital/arka_mid/models"
)

func TestValidarTrasladoInterno(t *testing.T) {
	tests := []struct {
		name    string
		detalle string
		wantErr string
	}{
		{
			name:    "detalle valido",
			detalle: `{"FuncionarioOrigen":10,"FuncionarioDestino":20,"Ubicacion":30,"Elementos":[1,2]}`,
		},
		{
			name:    "mismo funcionario",
			detalle: `{"FuncionarioOrigen":10,"FuncionarioDestino":10,"Ubicacion":30,"Elementos":[1,2]}`,
			wantErr: "mismo funcionario",
		},
		{
			name:    "origen invalido",
			detalle: `{"FuncionarioOrigen":0,"FuncionarioDestino":20,"Ubicacion":30,"Elementos":[1,2]}`,
			wantErr: "funcionario origen válido",
		},
		{
			name:    "destino invalido",
			detalle: `{"FuncionarioOrigen":10,"FuncionarioDestino":0,"Ubicacion":30,"Elementos":[1,2]}`,
			wantErr: "funcionario destino válido",
		},
		{
			name:    "detalle invalido",
			detalle: `{"FuncionarioOrigen":10,`,
			wantErr: "unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validarTrasladoInterno(&models.Movimiento{Detalle: tt.detalle})
			if tt.wantErr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}
