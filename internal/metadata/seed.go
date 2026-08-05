package metadata

import (
	"context"
	"database/sql"
	"fmt"
)

func Seed(ctx context.Context, db *sql.DB, cfg *Config) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, batch := range cfg.Batches {
		btc, err := upsertBatch(ctx, tx, batch)
		if err != nil {
			return err
		}
		fmt.Printf("Batch Created: %s(%s)\n", btc.BatchCode, btc.BatchID)

		for _, layer := range btc.Layers {
			lyr, err := upsertLayer(ctx, tx, layer)
			if err != nil {
				return err
			}
			fmt.Printf("Layer Created: %s(ID: %s). Part of Batch: %s(%s)\n", lyr.LayerCode, lyr.LayerID, btc.BatchCode, btc.BatchID)

			for _, module := range lyr.Modules {
				mdl, err := upsertModule(ctx, tx, module)
				if err != nil {
					return err
				}
				fmt.Printf("New module added: %s\n", mdl.ModuleID)

				for _, params := range mdl.Parameters {
					prm, err := UpsertParameters(ctx, tx, params)
					if err != nil {
						return err
					}
					fmt.Printf("New parameter added: %s\n", prm.ParameterID)

					if err := UpsertModuleParameters(ctx, tx,  mdl, prm); err != nil {
						return err
					}
				}

				if err := UpsertBatchModule(ctx, tx, btc, lyr, mdl); err != nil {
					return err
				}

				for _, params := range batch.Parameters {
					prm, err := UpsertParameters(ctx, tx, params)
					if err != nil {
						return err
					}

					if err := UpsertBatchParameters(ctx, tx, btc, prm); err != nil {
						return err
					}
				}
	 		}
		}
	}

	return tx.Commit()
}

func upsertBatch(ctx context.Context, tx *sql.Tx, batch Batch) (Batch, error) {
	batchUps := batch
	const query = `
  INSERT INTO dp_ctrl.BATCH
	(
		 BATCH_CODE
		,DESCRIPTION
		,IS_ACTIVE
		,CREATE_DATE
	)
  VALUES ($1, $2, $3, now())
	ON CONFLICT (BATCH_CODE)
	DO UPDATE
	SET
		DESCRIPTION = EXCLUDED.DESCRIPTION,
		IS_ACTIVE   = EXCLUDED.IS_ACTIVE
	RETURNING BATCH_ID;
  `

	err := tx.QueryRowContext(ctx, query, batch.BatchCode, batch.Description, batch.IsActive).Scan(&batchUps.BatchID)
	if err != nil {
		return Batch{}, err
	}

	return batchUps, nil
}

func upsertLayer(ctx context.Context, tx *sql.Tx, layer Layer) (Layer, error) {
	layerUps := layer
	const query = `
	INSERT INTO dp_ctrl.LAYER
	(
		LAYER_CODE
   ,DESCRIPTION
   ,IS_ACTIVE
   ,PRIORITY
   ,CREATE_DATE
	)
	VALUES ($1, $2, $3, $4, now())
	ON CONFLICT (LAYER_CODE)
	DO UPDATE
	SET
		DESCRIPTION = EXCLUDED.DESCRIPTION,
		IS_ACTIVE   = EXCLUDED.IS_ACTIVE
	RETURNING LAYER_ID
	`

	err := tx.QueryRowContext(ctx, query, layer.LayerCode, layer.Description, layer.IsActive, layer.Priority).Scan(&layerUps.LayerID)
	if err != nil {
		return Layer{}, err
	}

	return layerUps, nil
}

func upsertModule(ctx context.Context, tx *sql.Tx, module Module) (Module, error) {
	moduleUps := module
	const query = `
	INSERT INTO dp_ctrl.MODULE
	(
		MODULE_CODE
   ,DESCRIPTION
   ,IS_ACTIVE
   ,CREATE_DATE
	)
	VALUES ($1, $2, $3, now())
	ON CONFLICT (MODULE_CODE)
	DO UPDATE
	SET
		DESCRIPTION = EXCLUDED.DESCRIPTION,
		IS_ACTIVE   = EXCLUDED.IS_ACTIVE
	RETURNING MODULE_ID
	`

	err := tx.QueryRowContext(ctx, query, module.ModuleCode, module.Description, module.IsActive).Scan(&moduleUps.ModuleID)
	if err != nil {
		return Module{}, err
	}

	return moduleUps, nil
}

func UpsertParameters(ctx context.Context, tx *sql.Tx, params Parameter) (Parameter, error) {
	paramsUps := params
	const query = `
	INSERT INTO dp_ctrl.PARAMETER
	(
	  PARAMETER_CODE
	 ,DESCRIPTION
	 ,IS_ACTIVE
	 ,CREATE_DATE
	)
	VALUES ($1, $2, $3, now())
	ON CONFLICT (PARAMETER_CODE)
	DO UPDATE
	SET
		DESCRIPTION = EXCLUDED.DESCRIPTION,
		IS_ACTIVE   = EXCLUDED.IS_ACTIVE
	RETURNING PARAMETER_ID
	`
	err := tx.QueryRowContext(ctx, query, params.ParameterCode, params.ParameterCode, params.IsActive).Scan(&paramsUps.ParameterID)
	if err != nil {
		return Parameter{}, err
	}

	return paramsUps, nil
}

func UpsertBatchParameters(ctx context.Context, tx *sql.Tx, batch Batch, params Parameter) error {
	const query = `
	INSERT INTO dp_ctrl.BATCH_PARAMETERS
	(
	  BATCH_ID
	 ,PARAMETER_ID
	 ,PARAMETER_VALUE
	 ,CREATE_DATE
	)
	VALUES ($1, $2, $3, now())
	ON CONFLICT (BATCH_ID, PARAMETER_ID)
	DO UPDATE
	SET
		UPDATE_DATE = now()
	`
	_, err := tx.ExecContext(ctx, query, batch.BatchID, params.ParameterID, params.ParameterValue)
	if err != nil {
		return err
	}
	return nil
}

func UpsertModuleParameters(ctx context.Context, tx *sql.Tx, module Module, params Parameter) error {
	const query = `
	INSERT INTO dp_ctrl.MODULE_PARAMETERS
	(
	  MODULE_ID
	 ,PARAMETER_ID
	 ,PARAMETER_VALUE
	 ,CREATE_DATE
	)
	VALUES ($1, $2, $3, now())
	ON CONFLICT (MODULE_ID, PARAMETER_ID)
	DO UPDATE
	SET
		UPDATE_DATE = now()
	`
	_, err := tx.ExecContext(ctx, query, module.ModuleID, params.ParameterID, params.ParameterValue)
	if err != nil {
		return err
	}
	return nil
}

func UpsertBatchModule(ctx context.Context, tx *sql.Tx, batch Batch, layer Layer, module Module) error {
	const query = `
	INSERT INTO dp_ctrl.BATCH_MODULE
	(
	  BATCH_ID
	 ,MODULE_ID
	 ,MODULE_LAYER_ID
	 ,MODULE_SCRIPT
	 ,PRIORITY
	 ,IS_INITIAL_FULL
	 ,CREATE_DATE
	)
	VALUES ($1, $2, $3, $4, $5, $6, now())
	ON CONFLICT (BATCH_ID, MODULE_LAYER_ID, MODULE_ID)
	DO UPDATE
	SET
		UPDATE_DATE = now()
	`
	_, err := tx.ExecContext(ctx, query, batch.BatchID, module.ModuleID, layer.LayerID, module.Script, module.Priority, module.IsInitialFull)
	if err != nil {
		return err
	}

	return nil
}
