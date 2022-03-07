package main

// GXWrapTag is the way the texture is formatted
type GXWrapTag uint8

const (
	GX_CLAMP GXWrapTag = iota
	GX_REPEAT
	GX_MIRROR
	GX_MAXTEXWRAPMODE
)

type TexCoordGenTypes uint8

const (
	GX_TG_MTX3x4 TexCoordGenTypes = iota
	GX_TG_MTX2x4
	GX_TG_BUMP0
	GX_TG_BUMP1
	GX_TG_BUMP2
	GX_TG_BUMP3
	GX_TG_BUMP4
	GX_TG_BUMP5
	GX_TG_BUMP6
	GX_TG_BUMP7
	GX_TG_SRTG
)

type TexCoordGenSource uint8

const (
	GX_TG_POS TexCoordGenSource = iota
	GX_TG_NRM
	GX_TG_BINRM
	GX_TG_TANGENT
	GX_TG_TEX0
	GX_TG_TEX1
	GX_TG_TEX2
	GX_TG_TEX3
	GX_TG_TEX4
	GX_TG_TEX5
	GX_TG_TEX6
	GX_TG_TEX7
	GX_TG_TEXCOORD0
	GX_TG_TEXCOORD1
	GX_TG_TEXCOORD2
	GX_TG_TEXCOORD3
	GX_TG_TEXCOORD4
	GX_TG_TEXCOORD5
	GX_TG_TEXCOORD6
	GX_TG_COLOR0
	GX_TG_COLOR1
)

type TexCoordGenMatrixSource uint8

const (
	GX_PNMTX0 TexCoordGenMatrixSource = iota
	GX_PNMTX1
	GX_PNMTX2
	GX_PNMTX3
	GX_PNMTX4
	GX_PNMTX5
	GX_PNMTX6
	GX_PNMTX7
	GX_PNMTX8
	GX_PNMTX9
	GX_TEXMTX0
	GX_TEXMTX1
	GX_TEXMTX2
	GX_TEXMTX3
	GX_TEXMTX4
	GX_TEXMTX5
	GX_TEXMTX6
	GX_TEXMTX7
	GX_TEXMTX8
	GX_TEXMTX9
	GX_IDENTITY
	GX_DTTMTX0
	GX_DTTMTX1
	GX_DTTMTX2
	GX_DTTMTX3
	GX_DTTMTX4
	GX_DTTMTX5
	GX_DTTMTX6
	GX_DTTMTX7
	GX_DTTMTX8
	GX_DTTMTX9
	GX_DTTMTX10
	GX_DTTMTX11
	GX_DTTMTX12
	GX_DTTMTX13
	GX_DTTMTX14
	GX_DTTMTX15
	GX_DTTMTX16
	GX_DTTMTX17
	GX_DTTMTX18
	GX_DTTMTX19
	GX_DTTIDENTITY
)